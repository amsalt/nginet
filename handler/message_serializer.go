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
	ms := &MessageSerializer{}
	ms.DefaultOutboundHandler = core.NewDefaultOutboundHandler()
	ms.codec = codec
	ms.register = register
	return ms
}

func (ms *MessageSerializer) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	if rawBytes, ok := msg.([]byte); ok {
		ctx.FireWrite(rawBytes)
	} else if _, ok := msg.(bytes.WriteOnlyBuffer); ok {
		ctx.FireWrite(msg)
	} else if params, ok := msg.([]interface{}); ok && len(params) > 2 {
		idBuf := params[0]
		msgOrigin := params[1]
		msgID := params[2]
		buf, ok := idBuf.(bytes.WriteOnlyBuffer)
		if ok {
			err := ms.EncodePayload(buf, msgOrigin, msgID)
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

func (ms *MessageSerializer) EncodePayload(bufWithID bytes.WriteOnlyBuffer, msg interface{}, msgID interface{}) error {
	var rawData []byte
	var err error
	if raw, ok := msg.(*packet.RawPacket); ok {
		data, ok := raw.Payload().([]byte)
		if !ok {
			return errors.New("packet.RawPacket invalid payload.")
		}
		rawData = data

	} else {
		meta := ms.register.GetMetaByID(msgID)
		if meta != nil && meta.Codec() != nil {
			rawData, err = meta.Codec().Marshal(msg)
		} else {
			rawData, err = ms.codec.Marshal(msg)
		}

		if err != nil {
			log.Errorf("PacketSerializer.OnWrite Marshal failed for %+v", err)
			return err
		}
	}

	bufWithID.WriteTail(rawData)
	return nil
}
