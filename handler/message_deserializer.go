package handler

import (
	"errors"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/encoding"
	"github.com/amsalt/nginet/message"
)

type MessageDeserializer struct {
	*core.DefaultInboundHandler

	codec    encoding.Codec
	register message.Register
}

func NewMessageDeserializer(register message.Register, codec encoding.Codec) *MessageDeserializer {
	md := &MessageDeserializer{}
	md.DefaultInboundHandler = core.NewDefaultInboundHandler()
	md.codec = codec
	md.register = register
	return md
}

// OnRead ipdlements InboundHandler.
func (md *MessageDeserializer) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if params, ok := msg.([]interface{}); ok && len(params) > 1 {
		id := params[0]
		msgBuf, ok := params[1].(bytes.ReadOnlyBuffer)
		if ok {
			result, err := md.DecodePayload(id, msgBuf)

			if err == nil {
				var output []interface{}
				output = append(output, id, result)
				ctx.FireRead(output)
			} else {
				log.Errorf("PacketDeserializer.OnRead failed: %+v", err)
				ctx.FireError(err)
			}
		} else {
			ctx.FireError(errors.New("MessageDeserializer.OnRead invalid msg type, a bytes.ReadOnlyBuffer required."))
		}

	} else {
		ctx.FireError(errors.New("MessageDeserializer.OnRead invalid msg type, an array required."))
	}
}

func (md *MessageDeserializer) DecodePayload(msgID interface{}, data bytes.ReadOnlyBuffer) (interface{}, error) {
	meta := md.register.GetMetaByID(msgID)
	if meta != nil {
		msg := meta.CreateInstance()
		if meta.Codec() != nil { // support for meta specified codec.
			meta.Codec().Unmarshal(data.Bytes(), msg)
		} else {
			md.codec.Unmarshal(data.Bytes(), msg)
		}
		return msg, nil
	}
	return nil, errors.New("MessageDeserializer.DecodePayload message not registered")
}
