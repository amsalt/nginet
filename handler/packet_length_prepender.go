package handler

import (
	"encoding/binary"
	"errors"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
)

var ErrOutBoundEventType = errors.New("Invalid type for OnOutboundEvent args")

// transport layer message format, append size information to the head of Packet.
// 	size: segment size
//  id_header: msg id,defined by logic
//  msg_payload: msg content
//
// 					Segment
// --------------|--------------|--------------
// |   length    |  msg	serialized payload    |
// --------------|--------------|--------------
//
// binary data encoder will append packet length to the header.
// netty-like PacketLengthPrepender
type PacketLengthPrepender struct {
	*core.DefaultOutboundHandler
	lengthFieldLength uint
	byteorder         binary.ByteOrder // default binary.BigEndian
}

func NewPacketLengthPrepender(lengthFieldLength uint) *PacketLengthPrepender {
	plp := &PacketLengthPrepender{lengthFieldLength: lengthFieldLength}
	plp.DefaultOutboundHandler = core.NewDefaultOutboundHandler()
	plp.byteorder = binary.BigEndian

	return plp
}

// SetByteOrder Set byte order, default is binary.BigEndian
// 	byteorder:
// 		binary.BigEndian
// 		binary.LittleEndian
func (plp *PacketLengthPrepender) SetByteOrder(byteorder binary.ByteOrder) *PacketLengthPrepender {
	plp.byteorder = byteorder
	return plp
}

// OnOutboundEvent process outbound event.
func (plp *PacketLengthPrepender) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	if rawBytes, ok := msg.([]byte); ok {
		ctx.FireWrite(rawBytes)
	} else {
		data, err := plp.encode(ctx, msg)
		if err == nil {
			log.Debugf("PacketLengthPrepender.OnWrite result: %+v", data)
			ctx.FireWrite(data)
		} else {
			log.Errorf("PacketLengthPrepender.OnWrite failed: %+v", err)
			ctx.FireError(err)
		}
	}
}

// append fixed length field to header.
func (plp *PacketLengthPrepender) encode(ctx *core.ChannelContext, msg interface{}) (output interface{}, err error) {
	maxFrameLen := calcMaxFrameLen(plp.lengthFieldLength)
	if buff, ok := msg.(bytes.WriteOnlyBuffer); ok {
		actualLen := uint(buff.Len()) + plp.lengthFieldLength
		if uint64(actualLen) > maxFrameLen {
			return nil, ErrFrameTooLong
		}

		var head []byte
		switch plp.lengthFieldLength {
		case 1:
			head = []byte{byte(actualLen)}
			_, err = buff.WriteHeader(head)
		case 2:
			head = make([]byte, 2)
			plp.byteorder.PutUint16(head, uint16(actualLen))
			_, err = buff.WriteHeader(head)
		case 4:
			head = make([]byte, 4)
			plp.byteorder.PutUint32(head, uint32(actualLen))
			_, err = buff.WriteHeader(head)
		case 8:
			head = make([]byte, 8)
			plp.byteorder.PutUint64(head, uint64(actualLen))
			_, err = buff.WriteHeader(head)
		}

		output = buff
	} else {
		err = errors.New("PacketLengthPrepender msg should be bytes.WriteOnlyBuffer type")
	}

	return
}
