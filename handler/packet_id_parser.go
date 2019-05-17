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
	log.Debugf("IDParser OnRead message type: %T", msg)
	if data, ok := msg.(bytes.ReadOnlyBuffer); ok {
		output, err := ip.decodeId(data)
		if err == nil {
			ctx.FireRead(output)
		} else {
			ctx.FireError(err)
		}
	} else if data, ok := msg.([]interface{}); ok {
		log.Debugf("IDParser OnRead message len: %+v, data[0] type: %T", len(data), data[0])
		if len(data) >= 2 {
			if m, ok := data[0].(bytes.ReadOnlyBuffer); ok {
				output, err := ip.decodeId(m)
				if err == nil {
					output = append(output, data[1:]...)
					ctx.FireRead(output)
				} else {
					log.Errorf("IDParser err: %+v", err)
					ctx.FireError(err)
				}
			}
		}
	} else {
		log.Errorf("IDParser unsupported message type: %T", msg)
		ctx.FireError(fmt.Errorf("IDParser unsupported message type: %T", msg))
	}
}

func (ip *IDParser) decodeId(msg bytes.ReadOnlyBuffer) ([]interface{}, error) {
	msgID, msgBuf, err := ip.DecodeID(msg)
	log.Debugf("IDParser.OnRead parse msgID: %+v, msgBuf: %+v, err: %+v", msgID, msgBuf, err)
	if err == nil {
		var output []interface{}

		output = append(output, msgID, msgBuf)
		return output, nil
	} else {
		return nil, err
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
	log.Debugf("IDParser.OnWrite msg: %+v", msg)
	if rawBytes, ok := msg.([]byte); ok {
		ctx.FireWrite(rawBytes)
	} else if _, ok := msg.(bytes.WriteOnlyBuffer); ok {
		ctx.FireWrite(msg)
	} else {
		idBuf, id, err := ip.EncodeID(msg)
		log.Debugf("IDParser.OnWrite msg: %+v, id: %+v", msg, id)
		if err == nil {
			var output []interface{}
			output = append(output, idBuf, msg, id)
			ctx.FireWrite(output)
		} else {
			ctx.FireError(err)
		}
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

func (ip *IDParser) EncodeID(msg interface{}) (bytes.WriteOnlyBuffer, interface{}, error) {
	var msgID interface{}

	if raw, ok := msg.(*packet.RawPacket); ok {
		msgID = raw.ID()

	} else {
		meta := ip.register.GetMetaByMsg(msg)
		if meta == nil {
			return nil, nil, fmt.Errorf("IDParser.OnWrite msg: %+v not registered", msg)
		}
		msgID = meta.ID()
		log.Debugf("IDParser.EncodeID msg: %+v, id: %+v", msg, msgID)
	}

	// the max pre-reserved header size is MaxPacketLen + MaxExtraLen
	buf := bytes.NewWriteOnlyBuffer(MaxPacketLen + MaxExtraLen)
	err := ip.parser.Encode(msgID, buf)

	if err != nil {
		return nil, nil, fmt.Errorf("IDParser.OnWrite encode message id failed: %+v", err)

	}
	return buf, msgID, nil
}
