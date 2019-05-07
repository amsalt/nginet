package handler

import (
	"fmt"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/message"
	"github.com/amsalt/log"
)

// IDParser parses message ID.
// Deprecated
type IDParser struct {
	register message.Register
	parser   message.PacketIDParser
}

// NewIDParser creates new IDParser instance.
func NewIDParser(register message.Register, parser message.PacketIDParser) *IDParser {
	ip := &IDParser{}
	ip.parser = parser
	ip.register = register
	return ip
}

// OnRead implements InboundHandler
func (ip *IDParser) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if msg, ok := msg.(bytes.ReadOnlyBuffer); ok {
		var msgID interface{}
		msgID, err := ip.parser.Decode(msg)
		if err != nil {
			ctx.FireError(fmt.Errorf("IDParser.OnRead decode msg failed: %+v", err))
			return
		}

		var output []interface{}
		output = append(output, msgID)
		output = append(output, msg)
		ctx.FireRead(output)
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
	meta := ip.register.GetMetaByMsg(msg)
	if meta == nil {
		log.Errorf("IDParser.OnWrite msg: %+v not registered", msg)
		return
	}

	msgID := meta.ID()

	// the max pre-reserved header size is MaxPacketLen + MaxExtraLen
	buf := bytes.NewWriteOnlyBuffer(MaxPacketLen + MaxExtraLen)
	err := ip.parser.Encode(msgID, buf)

	if err != nil {
		log.Errorf("IDParser.OnWrite encode message id failed: %+v", err)
		return
	}

	log.Debugf("IDParser.OnWrite receive msg ------ msg: %+v\n", msg)
	var output []interface{}
	output = append(output, buf)
	output = append(output, msg)

	ctx.FireWrite(output)
	return

}
