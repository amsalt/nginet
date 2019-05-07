package handler

import (
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/message"
)

type DefaultMessageHandler struct {
	*core.DefaultInboundHandler
	processorMgr message.ProcessorMgr
}

func NewDefaultMessageHandler(mgr message.ProcessorMgr) *DefaultMessageHandler {
	mh := new(DefaultMessageHandler)
	mh.DefaultInboundHandler = core.NewDefaultInboundHandler()
	mh.processorMgr = mgr

	return mh
}

func (mh *DefaultMessageHandler) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if params, ok := msg.([]interface{}); ok && len(params) > 1 {
		id := params[0]
		p := mh.processorMgr.GetProcessorByID(id)
		p.SafeCall(ctx, params[1])
		ctx.FireRead(msg)
	}
}
