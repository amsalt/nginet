package core

// Handlers represents all handlers for kinds of network event(such as connect,
//  disconnect, inbound, outbound and error).

// InboundHandler processes all inbound event.
type InboundHandler interface {
	// OnConnect called when a new channel connected.
	OnConnect(ctx *ChannelContext, channel Channel)

	// OnDisconnect called when a channel disconnected.
	OnDisconnect(ctx *ChannelContext)

	// OnRead called when reads new data.
	OnRead(ctx *ChannelContext, msg interface{})

	// OnRead called when event triggered.
	OnEvent(ctx *ChannelContext, event interface{})

	// OnError called when error occurred
	OnError(ctx *ChannelContext, err error)
}

// OutboundHandler processes all outbound event.
type OutboundHandler interface {
	// OnWrite calls when write new data.
	OnWrite(ctx *ChannelContext, msg interface{})
}

type DefaultInboundHandler struct{}

func NewDefaultInboundHandler() *DefaultInboundHandler {
	return new(DefaultInboundHandler)
}

func (ih *DefaultInboundHandler) OnConnect(ctx *ChannelContext, channel Channel) {
	ctx.FireConnect(channel)
}

func (ih *DefaultInboundHandler) OnDisconnect(ctx *ChannelContext) {
	ctx.FireDisconnect()
}

func (ih *DefaultInboundHandler) OnRead(ctx *ChannelContext, msg interface{}) {
	ctx.FireRead(msg)
}

func (ih *DefaultInboundHandler) OnEvent(ctx *ChannelContext, event interface{}) {
	ctx.FireOnEvent(event)
}

func (ih *DefaultInboundHandler) OnError(ctx *ChannelContext, err error) {
	ctx.FireError(err)
}

type DefaultOutboundHandler struct{}

func NewDefaultOutboundHandler() *DefaultOutboundHandler {
	return new(DefaultOutboundHandler)
}

func (oh *DefaultOutboundHandler) OnWrite(ctx *ChannelContext, msg interface{}) {
	ctx.FireWrite(msg)
}
