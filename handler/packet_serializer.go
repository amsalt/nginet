package handler

import (
	"errors"
	"fmt"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/encoding"
	"github.com/amsalt/nginet/message"
	"github.com/amsalt/nginet/message/packet"
	"github.com/amsalt/log"
)

// PacketSerializer serializes packet.
// If input message is raw data, codec operation will be aborted.
type PacketSerializer struct {
	*core.DefaultOutboundHandler

	codec encoding.Codec

	register message.Register
	parser   message.PacketIDParser
}

// NewPacketSerializer creates new PacketSerializer instance.
func NewPacketSerializer(register message.Register, parser message.PacketIDParser, codec encoding.Codec) *PacketSerializer {
	ps := &PacketSerializer{}
	ps.DefaultOutboundHandler = core.NewDefaultOutboundHandler()
	ps.codec = codec
	ps.parser = parser
	ps.register = register
	return ps
}

func (ps *PacketSerializer) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	log.Debugf("PacketSerializer OnWrite: %+v", msg)
	buf, err := ps.encodeID(msg)
	if err != nil {
		ctx.FireError(err)
		return
	}

	err = ps.encodePayload(buf, msg)
	if err != nil {
		ctx.FireError(err)
		return
	}
	ctx.FireWrite(buf)
}

func (ps *PacketSerializer) encodePayload(bufWithID bytes.WriteOnlyBuffer, msg interface{}) error {
	var rawData []byte
	var err error
	if raw, ok := msg.(*packet.RawPacket); ok {
		data, ok := raw.Payload().([]byte)
		if !ok {
			return errors.New("packet.RawPacket invalid payload.")
		}

		rawData = data
	} else {
		rawData, err = ps.codec.Marshal(msg)
		if err != nil {
			log.Errorf("PacketSerializer.OnWrite Marshal failed for %+v", err)
			return err
		}
	}

	bufWithID.WriteTail(rawData)
	return nil
}

func (ps *PacketSerializer) encodeID(msg interface{}) (bytes.WriteOnlyBuffer, error) {
	var msgID interface{}

	if raw, ok := msg.(*packet.RawPacket); ok {
		msgID = raw.ID()

	} else {
		meta := ps.register.GetMetaByMsg(msg)
		if meta == nil {
			return nil, fmt.Errorf("IDParser.OnWrite msg: %+v not registered", msg)
		}
		msgID = meta.ID()
	}

	// the max pre-reserved header size is MaxPacketLen + MaxExtraLen
	buf := bytes.NewWriteOnlyBuffer(MaxPacketLen + MaxExtraLen)
	err := ps.parser.Encode(msgID, buf)

	if err != nil {
		return nil, fmt.Errorf("IDParser.OnWrite encode message id failed: %+v", err)

	}
	return buf, nil

}
