package handler

import (
	"github.com/amsalt/log"
	"github.com/amsalt/nginet/core"
)

// MessageEncoder serializes packet.
// If input message is raw data, codec operation will be aborted.
type MessageEncoder struct {
	*core.DefaultOutboundHandler

	idParser          *IDParser
	messageSerializer *MessageSerializer
}

// NewMessageEncoder creates new MessageEncoder instance.
func NewMessageEncoder(messageSerializer *MessageSerializer, idParser *IDParser) *MessageEncoder {
	me := &MessageEncoder{}
	me.DefaultOutboundHandler = core.NewDefaultOutboundHandler()
	me.idParser = idParser
	me.messageSerializer = messageSerializer
	return me
}

func (me *MessageEncoder) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	log.Debugf("MessageEncoder OnWrite: %+v", msg)
	if rawBytes, ok := msg.([]byte); ok {
		ctx.FireWrite(rawBytes)
	} else {
		buf, err := me.idParser.EncodeID(msg)
		if err != nil {
			ctx.FireError(err)
			return
		}

		err = me.messageSerializer.EncodePayload(buf, msg)
		if err != nil {
			ctx.FireError(err)
			return
		}
		ctx.FireWrite(buf)
	}
}
