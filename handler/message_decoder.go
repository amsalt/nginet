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
	pd := &MessageDecoder{}
	pd.DefaultInboundHandler = core.NewDefaultInboundHandler()
	pd.idParser = idParser
	pd.messageDeserializer = messageDeserializer
	return pd
}

// OnRead ipdlements InboundHandler.
func (pd *MessageDecoder) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if buf, ok := msg.(bytes.ReadOnlyBuffer); ok {
		id, msg, err := pd.idParser.DecodeID(buf)
		if err == nil {
			result, err := pd.messageDeserializer.DecodePayload(id, msg)

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
