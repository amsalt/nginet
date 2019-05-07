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
	ps := &MessageEncoder{}
	ps.DefaultOutboundHandler = core.NewDefaultOutboundHandler()
	ps.idParser = idParser
	ps.messageSerializer = messageSerializer
	return ps
}

func (ps *MessageEncoder) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	log.Debugf("MessageEncoder OnWrite: %+v", msg)
	buf, err := ps.idParser.EncodeID(msg)
	if err != nil {
		ctx.FireError(err)
		return
	}

	err = ps.messageSerializer.EncodePayload(buf, msg)
	if err != nil {
		ctx.FireError(err)
		return
	}
	ctx.FireWrite(buf)
}
