package core

import "github.com/amsalt/log"

// ChannelPipeline is a list of ChannelHandler which handles or intercepts inbound events and outbound operations of a
// ChannelPipeline implements an advanced form of the Intercepting Filter pattern
// to give a user full control over how an event is handled and how the ChannelHandler in a pipeline
// interact with each other.
// One ChannelPipeline per Channel.
type ChannelPipeline interface {
	InboundInvoker
	OutboundInvoker

	AddFirst(executor Executor, name string, h interface{})
	AddLast(executor Executor, name string, h interface{})
	AddAfter(afterName string, executor Executor, name string, handler interface{})
	AddBefore(beforeName string, executor Executor, name string, handler interface{})
	Channel() Channel
}

type channelPipeline struct {
	channel Channel

	head *HeadContext
	tail *TailContext
}

// NewChannelPipeline creates a new ChannelPipeline instance.
func NewChannelPipeline(channel Channel) ChannelPipeline {
	cp := &channelPipeline{channel: channel}
	cp.init()

	return cp
}

func (cp *channelPipeline) init() {
	cp.head = NewHeadContext(cp)
	cp.tail = NewTailContext(cp)

	cp.head.next = cp.tail.ChannelContext
	cp.tail.prev = cp.head.ChannelContext
}

func (cp *channelPipeline) Channel() Channel {
	return cp.channel
}

func (cp *channelPipeline) FireConnect(channel Channel) InboundInvoker {
	invokeConnect0(cp.head.ChannelContext, channel)
	return cp
}
func (cp *channelPipeline) FireDisconnect() InboundInvoker {
	invokeDisconnect0(cp.head.ChannelContext)
	return cp
}

func (cp *channelPipeline) FireError(err error) InboundInvoker {
	invokeError0(cp.head.ChannelContext, err)
	return cp
}

func (cp *channelPipeline) FireRead(msg interface{}) InboundInvoker {
	invokeRead0(cp.head.ChannelContext, msg)
	return cp
}

func (cp *channelPipeline) FireWrite(msg interface{}) OutboundInvoker {
	cp.tail.FireWrite(msg)
	return cp
}

func (cp *channelPipeline) AddFirst(executor Executor, name string, handler interface{}) {
	newCtx := NewDefaultChannelContext(executor, name, cp, handler)
	cp.addFirst0(newCtx.ChannelContext)
}

func (cp *channelPipeline) AddLast(executor Executor, name string, handler interface{}) {
	newCtx := NewDefaultChannelContext(executor, name, cp, handler)
	cp.addLast0(newCtx.ChannelContext)
}

func (cp *channelPipeline) AddAfter(afterName string, executor Executor, name string, handler interface{}) {
	newCtx := NewDefaultChannelContext(executor, name, cp, handler)
	cp.addAfter0(afterName, newCtx.ChannelContext)

	ctx := cp.head.ChannelContext
	for ctx != nil {
		if ctx.Name() == "HeadContext" {
			log.Debugf("+--------------------ChannelHandler list------------------------+")
		} else if ctx.Name() == "TailContext" {
			log.Debugf("+---------------------------------------------------------------+")
		} else {
			log.Debugf("+------ name:%-22s isInbound: %-4v isOutbound: %-4v", ctx.Name(), ctx.inbound, ctx.outbound)
		}
		ctx = ctx.next
	}
}

func (cp *channelPipeline) AddBefore(beforeName string, executor Executor, name string, handler interface{}) {
	newCtx := NewDefaultChannelContext(executor, name, cp, handler)
	cp.addBefore0(beforeName, newCtx.ChannelContext)

	ctx := cp.head.ChannelContext
	for ctx != nil {
		if ctx.Name() == "HeadContext" {
			log.Debugf("+--------------------ChannelHandler list------------------------+")
		} else if ctx.Name() == "TailContext" {
			log.Debugf("+---------------------------------------------------------------+")
		} else {
			log.Debugf("+------ name:%-22s isInbound: %-4v isOutbound: %-4v", ctx.Name(), ctx.inbound, ctx.outbound)
		}
		ctx = ctx.next
	}
}

func (cp *channelPipeline) addLast0(newCtx *ChannelContext) {
	prev := cp.tail.prev

	newCtx.prev = prev
	newCtx.next = cp.tail.ChannelContext

	prev.next = newCtx
	cp.tail.prev = newCtx
}

func (cp *channelPipeline) addFirst0(newCtx *ChannelContext) {
	next := cp.head.next

	newCtx.prev = cp.head.ChannelContext
	newCtx.next = next

	cp.head.next = newCtx
	next.prev = newCtx
}

func (cp *channelPipeline) addBefore0(name string, newCtx *ChannelContext) bool {
	search := cp.head.ChannelContext
	found := false
	for search != nil {
		if search.Name() == name {
			prev := search.prev
			newCtx.prev = prev
			newCtx.next = search

			prev.next = newCtx
			search.prev = newCtx

			found = true
			return found
		} else {
			search = search.next
		}
	}
	if !found {
		log.Warningf("ChannelPipline addBefore failed not found name: %+v", name)
	}
	return found
}

func (cp *channelPipeline) addAfter0(name string, newCtx *ChannelContext) bool {
	search := cp.head.ChannelContext
	found := false

	for search != cp.tail.ChannelContext {
		if search.Name() == name {
			next := search.next
			newCtx.prev = search
			newCtx.next = next

			search.next = newCtx
			next.prev = newCtx

			found = true
			return found
		} else {
			search = search.next
		}
	}

	if !found {
		log.Warningf("ChannelPipline addAfter failed not found name: %+v", name)
	}
	return found
}
