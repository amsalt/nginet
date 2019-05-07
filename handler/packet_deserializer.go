package handler

import (
	"errors"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/encoding"
	"github.com/amsalt/nginet/message"
	"github.com/amsalt/log"
)

// PacketDeserializer deserializers packet, and *io.Buffer is required for type of input message.
type PacketDeserializer struct {
	*core.DefaultInboundHandler

	codec    encoding.Codec
	register message.Register
	parser   message.PacketIDParser
}

// NewPacketDeserializer creates new PacketDeserializer instance.
func NewPacketDeserializer(register message.Register, parser message.PacketIDParser, codec encoding.Codec) *PacketDeserializer {
	pd := &PacketDeserializer{}
	pd.DefaultInboundHandler = core.NewDefaultInboundHandler()
	pd.codec = codec
	pd.parser = parser
	pd.register = register
	return pd
}

// OnRead ipdlements InboundHandler.
func (pd *PacketDeserializer) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if buf, ok := msg.(bytes.ReadOnlyBuffer); ok {
		id, msg, err := pd.decodeID(buf)
		if err == nil {
			result, err := pd.decodePayload(id, msg)

			if err == nil {
				var output []interface{}
				output = append(output, id, result)
				ctx.FireRead(output)
			} else {
				log.Errorf("PacketDeserializer.OnRead failed: %+v", err)
			}
		} else {
			log.Errorf("PacketDeserializer.OnRead failed: %+v", err)
		}
	}
}

func (pd *PacketDeserializer) decodePayload(msgID interface{}, data bytes.ReadOnlyBuffer) (interface{}, error) {
	meta := pd.register.GetMetaByID(msgID)
	if meta != nil {
		msg := meta.CreateInstance()
		if meta.Codec() != nil { // support for meta specified codec.
			meta.Codec().Unmarshal(data.Bytes(), msg)
		} else {
			pd.codec.Unmarshal(data.Bytes(), msg)
		}
		return msg, nil
	}
	return nil, errors.New("PacketDeserializer.decodePayload message not registered")
}

func (pd *PacketDeserializer) decodeID(msg bytes.ReadOnlyBuffer) (interface{}, bytes.ReadOnlyBuffer, error) {
	var msgID interface{}
	log.Debugf("PacketDeserializer.decodeID: %+v", msg)
	msgID, err := pd.parser.Decode(msg)
	if err != nil {
		log.Errorf("PacketDeserializer.decodeID decode msg failed: %+v", err)
		return nil, msg, err
	}

	return msgID, msg, nil
}
