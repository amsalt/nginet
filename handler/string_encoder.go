package handler

import (
	"fmt"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
)

type StringEncoder struct {
	*core.DefaultInboundHandler
	*core.DefaultOutboundHandler
}

func NewStringEncoder() *StringEncoder {
	sh := &StringEncoder{
		DefaultInboundHandler:  core.NewDefaultInboundHandler(),
		DefaultOutboundHandler: core.NewDefaultOutboundHandler(),
	}

	return sh
}

func (sc *StringEncoder) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if buff, ok := msg.(bytes.ReadOnlyBuffer); ok {
		bytes, err := buff.Read(0, buff.Len())
		if err != nil {
			ctx.FireError(fmt.Errorf("StringEncoder.OnRead msg failed: %+v", err))
		}
		ctx.FireRead(string(bytes))
	} else {
		ctx.FireError(fmt.Errorf("StringEncoder.OnRead msg not bytes.ReadOnlyBuffer"))
	}
}

func (sc *StringEncoder) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	if _, ok := msg.([]byte); ok {
		ctx.FireWrite(msg)
	} else if _, ok := msg.(bytes.WriteOnlyBuffer); ok {
		ctx.FireWrite(msg)
	} else if str, ok := msg.(string); ok {
		ctx.FireWrite([]byte(str))
	}
}
