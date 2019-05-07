package handler

import (
	"errors"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/encoding"
	"github.com/amsalt/nginet/message"
	"github.com/amsalt/nginet/message/packet"
)

type MessageSerializer struct {
	*core.DefaultOutboundHandler

	codec encoding.Codec

	register message.Register
}

// NewMessageSerializer creates new PacketSerializer instance.
func NewMessageSerializer(register message.Register, codec encoding.Codec) *MessageSerializer {
	ps := &MessageSerializer{}
	ps.DefaultOutboundHandler = core.NewDefaultOutboundHandler()
	ps.codec = codec
	ps.register = register
	return ps
}

func (ps *MessageSerializer) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	if params, ok := msg.([]interface{}); ok && len(params) > 1 {
		idBuf := params[0]
		msgOrigin := params[1]
		buf, ok := idBuf.(bytes.WriteOnlyBuffer)
		if ok {
			err := ps.EncodePayload(buf, msgOrigin)
			if err != nil {
				ctx.FireError(err)
				return
			}
			ctx.FireWrite(buf)
		} else {
			ctx.FireError(errors.New("MessageSerializer.OnWrite invalid msg type,an bytes.WriteOnlyBuffer required."))
		}

	} else {
		ctx.FireError(errors.New("MessageSerializer.OnWrite invalid msg type, an array required."))
	}

}

func (ps *MessageSerializer) EncodePayload(bufWithID bytes.WriteOnlyBuffer, msg interface{}) error {
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
