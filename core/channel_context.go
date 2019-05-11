package core

import (
	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
)

// ChannelContext enables a ChannelHandler to interact with its ChannelPipeline
// and other handlers. Among other things a handler can notify the next ChannelHandler in the
// ChannelPipeline.
// a ChannelContext is an InboundInvoker or an OutboundInvoker or both.
type ChannelContext struct {
	next *ChannelContext
	prev *ChannelContext

	pipeline ChannelPipeline
	channel  Channel

	// to reduce type assertion when invoke.
	inHandler  InboundHandler
	outHandler OutboundHandler

	inbound  bool
	outbound bool

	name     string
	executor Executor
}

// newChannelContext creates a new ChannelContext
func newChannelContext(executor Executor, name string, pipeline ChannelPipeline, handler interface{}) *ChannelContext {
	ctx := &ChannelContext{name: name, executor: executor}
	ctx.init(pipeline, handler)
	if name == "HeadContext" {
		log.Debugf("+--------------------ChannelHandler list------------------------+")
	} else if name == "TailContext" {
		log.Debugf("+---------------------------------------------------------------+")
	} else {
		log.Debugf("+------ name:%-22s isInbound: %-4v isOutbound: %-4v", name, ctx.inbound, ctx.outbound)
	}
	return ctx
}
func (ctx *ChannelContext) init(pipeline ChannelPipeline, handler interface{}) {
	ctx.pipeline = pipeline

	in, ok := handler.(InboundHandler)
	ctx.inbound = ok
	if ok {
		ctx.inHandler = in
	}

	out, ok := handler.(OutboundHandler)
	ctx.outbound = ok
	if ok {
		ctx.outHandler = out
	}
}

func (ctx *ChannelContext) findNextInboundContext() *ChannelContext {
	var next *ChannelContext
	next = ctx.next
	for !next.inbound { // the last one is TailContext.
		next = next.next
	}

	return next
}

func (ctx *ChannelContext) findNextOutboundContext() *ChannelContext {
	var prev *ChannelContext
	prev = ctx.prev
	for !prev.outbound { // the last one is HeadContext.
		prev = prev.prev
	}

	return prev
}

func (ctx *ChannelContext) Channel() Channel {
	return ctx.pipeline.Channel()
}

func (ctx *ChannelContext) Attr() *AttrMap {
	return ctx.Channel().Attr()
}

func (ctx *ChannelContext) Name() string {
	return ctx.name
}

func (ctx *ChannelContext) Executor() Executor {
	return ctx.executor
}

func (ctx *ChannelContext) FireConnect(channel Channel) InboundInvoker {
	invokeConnect0(ctx.findNextInboundContext(), channel)
	return ctx
}

func invokeConnect0(next *ChannelContext, channel Channel) {
	if next.Executor() != nil {
		next.Executor().Execute(func() { next.doConnect(channel) })
	} else {
		next.doConnect(channel)
	}
}

func (ctx *ChannelContext) doConnect(channel Channel) {
	ctx.inHandler.OnConnect(ctx, channel)
}

func (ctx *ChannelContext) FireError(err error) InboundInvoker {
	invokeError0(ctx.findNextInboundContext(), err)
	return ctx
}

func invokeError0(next *ChannelContext, err error) {
	if next.Executor() != nil {
		next.Executor().Execute(func() { next.doError(err) })
	} else {
		next.doError(err)
	}
}

func (ctx *ChannelContext) doError(err error) {
	ctx.inHandler.OnError(ctx, err)
}

func (ctx *ChannelContext) FireDisconnect() InboundInvoker {
	invokeDisconnect0(ctx.findNextInboundContext())
	return ctx
}

func invokeDisconnect0(next *ChannelContext) {
	if next.Executor() != nil {
		next.Executor().Execute(func() {
			next.doDisconnect()
		})
	} else {
		next.doDisconnect()
	}

}

func (ctx *ChannelContext) doDisconnect() {
	ctx.inHandler.OnDisconnect(ctx)
}

func (ctx *ChannelContext) FireRead(msg interface{}) InboundInvoker {
	invokeRead0(ctx.findNextInboundContext(), msg)
	return ctx
}

func invokeRead0(next *ChannelContext, msg interface{}) {
	if next.Executor() != nil {
		next.Executor().Execute(func() {
			next.doRead(msg)
		})
	} else {
		next.doRead(msg)
	}

}

func (ctx *ChannelContext) doRead(msg interface{}) {
	ctx.inHandler.OnRead(ctx, msg)
}

func (ctx *ChannelContext) Write(msg interface{}) OutboundInvoker {
	ctx.Channel().Write(msg)
	return ctx
}

func (ctx *ChannelContext) Close() OutboundInvoker {
	ctx.Channel().Close()
	return ctx
}

func (ctx *ChannelContext) FireWrite(msg interface{}) OutboundInvoker {
	invokeWrite0(ctx.findNextOutboundContext(), msg)
	return ctx
}

func invokeWrite0(next *ChannelContext, msg interface{}) {
	log.Debugf("invokeWrite0 execute handler: %+v", next.Name())
	if next.Executor() != nil {
		// run in the special executor all that called after this handler, include this one.
		next.Executor().Execute(func() {
			next.doWrite(msg)
		})
	} else {
		next.doWrite(msg)
	}

}

func (ctx *ChannelContext) doWrite(msg interface{}) {
	ctx.outHandler.OnWrite(ctx, msg)
}

type DefaultChannelContext struct {
	*ChannelContext

	handler interface{}
}

func NewDefaultChannelContext(executor Executor, name string, pipeline ChannelPipeline, handler interface{}) *DefaultChannelContext {
	ctx := newChannelContext(executor, name, pipeline, handler)
	dctx := &DefaultChannelContext{}
	dctx.handler = handler
	dctx.ChannelContext = ctx

	return dctx
}

// HeadContext the header of pipeline.
type HeadContext struct {
	*ChannelContext
}

// NewHeadContext return a new instance of *HeadContext.
func NewHeadContext(pipeline ChannelPipeline) *HeadContext {
	hctx := &HeadContext{}
	hctx.ChannelContext = newChannelContext(nil, "HeadContext", pipeline, hctx)
	return hctx
}

// OnRead processes a read event.
func (hctx *HeadContext) OnRead(ctx *ChannelContext, msg interface{}) {
	ctx.FireRead(msg)
}

func (hctx *HeadContext) OnError(ctx *ChannelContext, err error) {
	log.Debugf("HeadContext OnError: %+v", err)
	ctx.FireError(err)
}

// OnConnect processes a connect event.
func (hctx *HeadContext) OnConnect(ctx *ChannelContext, channel Channel) {
	hctx.FireConnect(channel)
}

// OnDisconnect processes a disconnect event.
func (hctx *HeadContext) OnDisconnect(ctx *ChannelContext) {
	ctx.FireDisconnect()
}

// OnWrite processes a write event.
func (hctx *HeadContext) OnWrite(ctx *ChannelContext, msg interface{}) {
	if msg == nil {
		return
	}

	data, ok := msg.([]byte)
	if ok {
		ctx.Channel().RawConn().Write(data)
	} else if data, ok := msg.(bytes.WriteOnlyBuffer); ok {
		ctx.Channel().RawConn().Write(data.Bytes())
	} else {
		log.Errorf("HeadContext.OnWrite write with unsupported type: %T", msg)
	}
}

// TailContext represents the tail of pipeline.
type TailContext struct {
	*ChannelContext
}

// NewTailContext returns a new instance of *TailContext.
func NewTailContext(pipeline ChannelPipeline) *TailContext {
	tctx := &TailContext{}
	tctx.ChannelContext = newChannelContext(nil, "TailContext", pipeline, tctx)
	return tctx
}

// OnRead do nothing to stop the pipeline.
func (tctx *TailContext) OnRead(ctx *ChannelContext, msg interface{}) { /*do nothing: stop at tail*/ }

// OnConnect do nothing to stop the pipeline.
func (tctx *TailContext) OnConnect(ctx *ChannelContext, channel Channel) { /*do nothing: stop at tail*/ }

// OnDisconnect do nothing to stop the pipeline.
func (tctx *TailContext) OnDisconnect(ctx *ChannelContext) { /*do nothing: stop at tail*/ }

func (tctx *TailContext) OnError(ctx *ChannelContext, err error) { /*do nothing: stop at tail*/ }
