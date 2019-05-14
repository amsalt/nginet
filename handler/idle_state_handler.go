package handler

import (
	"time"

	"github.com/amsalt/nginet/core"
)

const (
	ReadTimeout  = iota
	WriteTimeout = iota
	AllTimeout   = iota
)

type IdleEvent struct {
	TimeoutType int
}

type IdleStateHandler struct {
	*core.DefaultInboundHandler
	*core.DefaultOutboundHandler
	readTimeout    int
	writeTimeout   int
	needAllTimeout bool
	rTimeout       bool
	wTimeout       bool

	lastReadTime      int
	lastWriteTime     int
	readTimeoutDelay  int
	writeTimeoutDelay int

	readCheckTimer  *time.Timer
	writeCheckTimer *time.Timer

	stop bool
}

func NewIdleStateHandler(readTimeoutSec, writeTimeoutSec int, needAllTimeout bool) *IdleStateHandler {
	ish := &IdleStateHandler{
		DefaultInboundHandler:  core.NewDefaultInboundHandler(),
		DefaultOutboundHandler: core.NewDefaultOutboundHandler(),
	}

	ish.readTimeout = readTimeoutSec
	ish.writeTimeout = writeTimeoutSec
	ish.needAllTimeout = needAllTimeout

	ish.readTimeoutDelay = readTimeoutSec
	ish.writeTimeoutDelay = writeTimeoutSec

	return ish
}

func (ish *IdleStateHandler) OnConnect(ctx *core.ChannelContext, channel core.Channel) {
	ish.run(ctx)
	ctx.FireConnect(channel)
}

func (ish *IdleStateHandler) OnDisconnect(ctx *core.ChannelContext) {
	ish.stop = true
	ish.abort(ctx)
	ctx.FireDisconnect()
}

func (ish *IdleStateHandler) OnRead(ctx *core.ChannelContext, msg interface{}) {
	ish.lastReadTime = time.Now().Nanosecond()
	ish.readTimeoutDelay = ish.readTimeout
	ctx.FireRead(msg)
}

func (ish *IdleStateHandler) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	ish.lastWriteTime = time.Now().Nanosecond()
	ish.writeTimeoutDelay = ish.writeTimeout
	ctx.FireWrite(msg)
}

func (ish *IdleStateHandler) channelIdle(ctx *core.ChannelContext, event *IdleEvent) {
	ctx.FireEvent(event)
}

func (ish *IdleStateHandler) run(ctx *core.ChannelContext) {
	ish.checkReadTimeout(ctx)
	ish.checkWriteTimeout(ctx)
}

func (ish *IdleStateHandler) abort(ctx *core.ChannelContext) {
	ish.readCheckTimer.Stop()
	ish.writeCheckTimer.Stop()
}

func (ish *IdleStateHandler) checkReadTimeout(ctx *core.ChannelContext) {
	if ish.stop {
		return
	}

	ish.readCheckTimer = time.AfterFunc(time.Second*time.Duration(ish.readTimeoutDelay), func() {
		ish.readTimeoutDelay -= (time.Now().Nanosecond() - ish.lastReadTime)
		if ish.readTimeoutDelay < 0 {
			ish.rTimeout = true
			if !ish.needAllTimeout {
				ish.channelIdle(ctx, &IdleEvent{TimeoutType: ReadTimeout})
				ish.rTimeout = false
				ish.readTimeoutDelay = ish.readTimeout
			} else if ish.wTimeout {
				ish.channelIdle(ctx, &IdleEvent{TimeoutType: AllTimeout})
				ish.rTimeout = false
				ish.wTimeout = false
				ish.readTimeoutDelay = ish.readTimeout
				ish.writeTimeoutDelay = ish.writeTimeout
			}
		}

		ish.checkReadTimeout(ctx)
	})
}

func (ish *IdleStateHandler) checkWriteTimeout(ctx *core.ChannelContext) {
	if ish.stop {
		return
	}

	ish.writeCheckTimer = time.AfterFunc(time.Second*time.Duration(ish.writeTimeoutDelay), func() {
		ish.writeTimeoutDelay -= (time.Now().Nanosecond() - ish.lastWriteTime)
		if ish.writeTimeoutDelay < 0 {
			ish.rTimeout = true
			if !ish.needAllTimeout {
				ish.channelIdle(ctx, &IdleEvent{TimeoutType: WriteTimeout})
				ish.wTimeout = false
				ish.writeTimeoutDelay = ish.writeTimeout
			} else if ish.rTimeout {
				ish.channelIdle(ctx, &IdleEvent{TimeoutType: AllTimeout})
				ish.rTimeout = false
				ish.wTimeout = false
				ish.readTimeoutDelay = ish.readTimeout
				ish.writeTimeoutDelay = ish.writeTimeout
			}
		}

		ish.checkWriteTimeout(ctx)
	})
}
