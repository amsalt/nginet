package handler

import (
	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
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
	if _, ok := msg.([]byte); ok {
		ctx.FireWrite(msg)
	} else if _, ok := msg.(bytes.WriteOnlyBuffer); ok {
		ctx.FireWrite(msg)
	} else {
		buf, err := me.encode(msg)
		if err != nil {
			ctx.FireError(err)
		} else {
			ctx.FireWrite(buf)
		}

	}
}

func (me *MessageEncoder) encode(msg interface{}) (bytes.WriteOnlyBuffer, error) {
	buf, id, err := me.idParser.EncodeID(msg)
	if err != nil {
		return nil, err
	}

	err = me.messageSerializer.EncodePayload(buf, msg, id)
	return buf, err
}
