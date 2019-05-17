package handler

import (
	"errors"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/message"
)

// DefaultMessageHandler represents a default implementation of MessageHandler.
type DefaultMessageHandler struct {
	*core.DefaultInboundHandler
	processorMgr message.ProcessorMgr
}

// NewDefaultMessageHandler creates and return a pointer to the instance of DefaultMessageHandler
func NewDefaultMessageHandler(mgr message.ProcessorMgr) *DefaultMessageHandler {
	mh := new(DefaultMessageHandler)
	mh.DefaultInboundHandler = core.NewDefaultInboundHandler()
	mh.processorMgr = mgr

	return mh
}

// OnRead InboundHandler
func (mh *DefaultMessageHandler) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if params, ok := msg.([]interface{}); ok && len(params) > 1 {
		id := params[0]
		p := mh.processorMgr.GetProcessorByID(id)
		if p == nil {
			log.Errorf("msg id %+v not register processor", id)
			ctx.FireError(errors.New("msg not registered"))
		}
		if len(params) > 2 {
			var args []interface{}
			args = append(args, params[2:]...)
			p.SafeCall(ctx, params[1], args...)
		} else {
			p.SafeCall(ctx, params[1])
		}

		ctx.FireRead(msg)
	} else {
		ctx.FireError(errors.New("DefaultMessageHandler.OnRead invalid msg format"))
	}
}
