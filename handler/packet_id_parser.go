package handler

import (
	"fmt"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/message"
	"github.com/amsalt/nginet/message/packet"
)

// IDParser parses message ID.
// Deprecated
type IDParser struct {
	*core.DefaultInboundHandler
	*core.DefaultOutboundHandler

	register message.Register
	parser   message.PacketIDParser
}

// NewIDParser creates new IDParser instance.
func NewIDParser(register message.Register, parser message.PacketIDParser) *IDParser {
	ip := &IDParser{}
	ip.parser = parser
	ip.register = register
	ip.DefaultInboundHandler = core.NewDefaultInboundHandler()
	ip.DefaultOutboundHandler = core.NewDefaultOutboundHandler()
	return ip
}

// OnRead implements InboundHandler
func (ip *IDParser) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if msg, ok := msg.(bytes.ReadOnlyBuffer); ok {
		msgID, msgBuf, err := ip.DecodeID(msg)
		if err == nil {
			var output []interface{}
			output = append(output, msgID, msgBuf)
			ctx.FireRead(output)
		} else {
			ctx.FireError(err)
		}

	}

}

func (ip *IDParser) OnConnect(ctx *core.ChannelContext, channel core.Channel) {
	ctx.FireConnect(channel)
}

// OnDisconnect called when channel disconnected.
func (ip *IDParser) OnDisconnect(ctx *core.ChannelContext) {
	ctx.FireDisconnect()
}

// OnWrite implements OutboundHandler
// outputs:
// 	an array arr
// 		arr[0] is bytes.WriteOnlyBuffer with message ID.
// 		arr[1] is original message object.
func (ip *IDParser) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	idBuf, err := ip.EncodeID(msg)
	if err == nil {
		var output []interface{}
		output = append(output, idBuf, msg)
		ctx.FireWrite(output)
	} else {
		ctx.FireError(err)
	}
}

func (ip *IDParser) DecodeID(msg bytes.ReadOnlyBuffer) (interface{}, bytes.ReadOnlyBuffer, error) {
	var msgID interface{}
	log.Debugf("PacketDeserializer.decodeID: %+v", msg)
	msgID, err := ip.parser.Decode(msg)
	if err != nil {
		log.Errorf("PacketDeserializer.decodeID decode msg failed: %+v", err)
		return nil, msg, err
	}

	return msgID, msg, nil
}

func (ip *IDParser) EncodeID(msg interface{}) (bytes.WriteOnlyBuffer, error) {
	var msgID interface{}

	if raw, ok := msg.(*packet.RawPacket); ok {
		msgID = raw.ID()

	} else {
		meta := ip.register.GetMetaByMsg(msg)
		if meta == nil {
			return nil, fmt.Errorf("IDParser.OnWrite msg: %+v not registered", msg)
		}
		msgID = meta.ID()
	}

	// the max pre-reserved header size is MaxPacketLen + MaxExtraLen
	buf := bytes.NewWriteOnlyBuffer(MaxPacketLen + MaxExtraLen)
	err := ip.parser.Encode(msgID, buf)

	if err != nil {
		return nil, fmt.Errorf("IDParser.OnWrite encode message id failed: %+v", err)

	}
	return buf, nil
}
