package handler

import (
	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
)

// MessageDecoder deserializers packet, and *io.Buffer is required for type of input message.
type MessageDecoder struct {
	*core.DefaultInboundHandler

	idParser            *IDParser
	messageDeserializer *MessageDeserializer
}

// NewMessageDecoder creates new MessageDecoder instance.
func NewMessageDecoder(messageDeserializer *MessageDeserializer, idParser *IDParser) *MessageDecoder {
	md := &MessageDecoder{}
	md.DefaultInboundHandler = core.NewDefaultInboundHandler()
	md.idParser = idParser
	md.messageDeserializer = messageDeserializer
	return md
}

// OnRead ipdlements InboundHandler.
func (md *MessageDecoder) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if buf, ok := msg.(bytes.ReadOnlyBuffer); ok {
		id, msg, err := md.idParser.DecodeID(buf)
		if err == nil {
			result, err := md.messageDeserializer.DecodePayload(id, msg)

			if err == nil {
				var output []interface{}
				output = append(output, id, result)
				ctx.FireRead(output)
			} else {
				log.Errorf("MessageDecoder.OnRead failed: %+v", err)
			}
		} else {
			log.Errorf("MessageDecoder.OnRead failed: %+v", err)
		}
	}
}
