package handler

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
)

const (
	MaxPacketLen = 8
	MaxExtraLen  = 8
)

var (
	ErrFrameTooLong     = errors.New("LengthFieldBasedFrameDecoder: frame too long")
	ErrInBoundEventType = errors.New("Invalid type for OnInboundEvent args")
	ErrInBoundMsgID     = errors.New("Invalid id for OnInboundEvent msg")
	ErrInBoundRPCSeq    = errors.New("Invalid rpc seq for OnInboundEvent msg")
	ErrNoMetaFound      = errors.New("No meta found")
)

// transport layer message format, append size information to the head of Packet.
// 	size: segment size
//  id_header: msg id,defined by logic
//  msg_payload: msg content
//
// 					Segment
// --------------|--------------|--------------
// |   length    |  id_header   | msg_payload |
// --------------|--------------|--------------
//
// binary data decoder based on head size and payload body
// netty-like PacketLengthDecoder
type PacketLengthDecoder struct {
	*core.DefaultInboundHandler

	maxFrameLength      uint64
	lengthFieldOffset   uint
	sizeOfLengthField   uint
	initialBytesToStrip int
	byteorder           binary.ByteOrder // default binary.BigEndian

	lengthFieldEndOffset int
}

// NewPacketLengthDecoder return new LengthFieldBasedFrameDecoder
// 		maxFrameLength: max size of msg
// 		lengthFieldOffset: the location of "lenght field" need offset
// 		lengthFieldLength: the length of "lenght field"
//		initialBytesToStrip: how many bytes need to be strip before delivering to next processor
//
func NewPacketLengthDecoder(sizeOfLengthField uint) *PacketLengthDecoder {
	pld := &PacketLengthDecoder{
		sizeOfLengthField: sizeOfLengthField,
	}
	pld.DefaultInboundHandler = core.NewDefaultInboundHandler()
	pld.maxFrameLength = calcMaxFrameLen(sizeOfLengthField) // default length
	pld.recalcLengthFieldEndOffset()
	pld.byteorder = binary.BigEndian
	pld.initialBytesToStrip = -1

	return pld
}

// SetByteOrder Set byte order, default is binary.BigEndian
// 	byteorder:
// 		binary.BigEndian
// 		binary.LittleEndian
func (pld *PacketLengthDecoder) SetByteOrder(byteorder binary.ByteOrder) *PacketLengthDecoder {
	pld.byteorder = byteorder
	return pld
}

// SetMaxFrameLength set the max value of package length.
func (pld *PacketLengthDecoder) SetMaxFrameLength(maxFrameLength uint64) *PacketLengthDecoder {
	// ignore if bigger then max value
	if maxFrameLength > pld.maxFrameLength {
		return pld
	}

	pld.maxFrameLength = maxFrameLength
	return pld
}

// SetLengthFieldOffset set the offset of length field,
// previous value maybe version or other value.
func (pld *PacketLengthDecoder) SetLengthFieldOffset(lengthFieldOffset uint) *PacketLengthDecoder {
	pld.lengthFieldOffset = lengthFieldOffset
	pld.recalcLengthFieldEndOffset()

	return pld
}

// SetInitialBytesToStrip set the number to skip when pass to next processor.
func (pld *PacketLengthDecoder) SetInitialBytesToStrip(initialBytesToStrip uint) *PacketLengthDecoder {
	pld.initialBytesToStrip = int(initialBytesToStrip)
	return pld
}

// OnRead implements OutboundHandler.
func (pld *PacketLengthDecoder) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if data, ok := msg.(bytes.ReadOnlyBuffer); ok {
		output, err := pld.decode(ctx, data)
		log.Debugf("PacketLengthDecoder.OnRead: %+v, err: %+v", output, err)
		if err == nil {
			ctx.FireRead(output)
		} else {
			ctx.FireError(fmt.Errorf("PacketLengthDecoder failed: %+v", err))
		}
	} else {
		ctx.FireError(fmt.Errorf("PacketLengthDecoder.OnRead msg not bytes.ReadOnlyBuffer"))
	}
}

func (pld *PacketLengthDecoder) decode(ctx *core.ChannelContext, buff bytes.ReadOnlyBuffer) (interface{}, error) {
	if buff.Len() < pld.lengthFieldEndOffset {
		return nil, errors.New("return for buff len less then lengthFieldEndOffset")
	}

	frameLength := pld.getUnadjustedFrameLength(buff.Bytes())
	log.Debugf("PacketLengthDecoder.decode frameLength: %+v", frameLength)

	if buff.Len() < int(frameLength) {
		errInfo := fmt.Sprintf("return for buff len: %v less then int(frameLength): %v ",
			buff.Len(), int(frameLength))
		log.Errorf(errInfo)
		return nil, errors.New(errInfo)
	}

	if frameLength > pld.maxFrameLength {
		return nil, errors.New("return for ErrFrameTooLong")
	}

	skip := pld.getInitialBytesToStrip()
	actualFrameLen := frameLength - skip

	// TODO: check buffer safe again.
	// Notice: buff must used in the same goroutine.
	msg, err := buff.Read(int(skip), int(actualFrameLen))
	if err != nil {
		errInfo := fmt.Sprintf("PacketLengthDecoder.decode buff read error: %+v", err)
		log.Errorf(errInfo)
		return nil, errors.New(errInfo)
	}

	log.Debugf("PacketLengthDecoder.decode: %+v", msg)
	return bytes.NewReadOnlyBufferWithBytes(msg), nil
}

func (pld *PacketLengthDecoder) getUnadjustedFrameLength(buf []byte) uint64 {
	var msgLen uint64
	switch pld.sizeOfLengthField {
	case 1:
		msgLen = uint64(buf[pld.lengthFieldOffset:pld.lengthFieldEndOffset][0])
	case 2:
		msgLen = uint64(pld.byteorder.Uint16(buf[pld.lengthFieldOffset:pld.lengthFieldEndOffset]))
	case 4:
		msgLen = uint64(pld.byteorder.Uint32(buf[pld.lengthFieldOffset:pld.lengthFieldEndOffset]))
	case 8:
		msgLen = uint64(pld.byteorder.Uint64(buf[pld.lengthFieldOffset:pld.lengthFieldEndOffset]))
	}
	return msgLen
}

func (pld *PacketLengthDecoder) recalcLengthFieldEndOffset() {
	pld.lengthFieldEndOffset = int(pld.lengthFieldOffset + pld.sizeOfLengthField)
}

// default skip lengthFieldOffset & sizeOfLengthField
// if don't want to skip, use SetInitialBytesToStrip(0)
func (pld *PacketLengthDecoder) getInitialBytesToStrip() uint64 {
	if pld.initialBytesToStrip == -1 {
		return uint64(pld.lengthFieldOffset + pld.sizeOfLengthField)
	}
	return uint64(pld.initialBytesToStrip)
}

func calcMaxFrameLen(lengthFieldLength uint) uint64 {
	var maxMsgLen uint64
	switch lengthFieldLength {
	case 1:
		maxMsgLen = 1<<8 - 1
	case 2:
		maxMsgLen = 1<<16 - 1
	case 4:
		maxMsgLen = 1<<32 - 1
	case 8:
		maxMsgLen = 1<<64 - 1
	}
	return maxMsgLen
}
