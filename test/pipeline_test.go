package test

import (
	"testing"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
)

type Msg struct {
	Hello string
}

type inhandler1 struct {
}

func (ih *inhandler1) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if m, ok := msg.(bytes.ReadOnlyBuffer); ok {
		log.Infof("inhandler1 msg: %+v", string(m.Bytes()))
	} else {
		log.Infof("inhandler1 msg: %+v", msg)
	}

	ctx.FireRead(msg)

	// var head = make([]byte, 2)
	// binary.BigEndian.PutUint16(head, 7)
	// time.Sleep(time.Second)
	// ctx.Write(bytes.NewWriteOnlyBufferWithBytes(8, []byte{head[0], head[1], 'h', 'e', 'l', 'l', 'o'}))
}

// OnConnect called when new channel connected.
func (ih *inhandler1) OnConnect(ctx *core.ChannelContext, channel core.Channel) {
	log.Infof("new client connected: %+v", ctx)
	ctx.FireConnect(channel)
}

func (ih *inhandler1) OnError(ctx *core.ChannelContext, err error) {
	log.Infof("OnError: %+v", err)
	ctx.FireError(err)
}

func (ih *inhandler1) OnEvent(ctx *core.ChannelContext, event interface{}) {
	log.Infof("OnEvent: %+v", event)
	ctx.FireEvent(event)
	ctx.Close()
}

// OnDisconnect called when channel disconnected.
func (ih *inhandler1) OnDisconnect(ctx *core.ChannelContext) {
	log.Infof("client disconnected: %+v", ctx)
	ctx.FireDisconnect()
}

type inhandler2 struct {
}

func (ih *inhandler2) OnRead(ctx *core.ChannelContext, msg interface{}) {
	buf := msg.(bytes.ReadOnlyBuffer)

	log.Debugf("inhandler2 read data: %+v", string(buf.Bytes()))
	ctx.FireRead(msg)
	// time.Sleep(time.Second)
	// ctx.Write([]byte("client response data"))
}

// OnConnect called when new channel connected.
func (ih *inhandler2) OnConnect(ctx *core.ChannelContext, channel core.Channel) {
	ctx.FireConnect(channel)
}

// OnDisconnect called when channel disconnected.
func (ih *inhandler2) OnDisconnect(ctx *core.ChannelContext) {
	log.Infof("server disconnected: %+v", ctx)
	ctx.FireDisconnect()
}

func (ih *inhandler2) OnEvent(ctx *core.ChannelContext, event interface{}) {
	log.Infof("OnEvent: %+v", event)
	ctx.FireEvent(event)
}

func (ih *inhandler2) OnError(ctx *core.ChannelContext, err error) {
	log.Infof("OnError: %+v", err)
	ctx.FireError(err)
}

type outhandler struct {
}

func (oh *outhandler) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	log.Infof("outhandler.OnWrite msg: %+v", msg)
	if data, ok := msg.(bytes.WriteOnlyBuffer); ok {
		log.Infof("outhandler.OnWrite bytes.WriteOnlyBuffer bytes: %+v", data.Bytes())
	}
	ctx.FireWrite(msg)
}

type outhandler2 struct {
}

func (oh *outhandler2) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	log.Infof("outhandler2 write msg: %+v", msg)
	ctx.FireWrite(msg)
}

type inouthandler struct{}

func (ih *inouthandler) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	ctx.FireWrite(msg)
}

func (ih *inouthandler) OnRead(ctx *core.ChannelContext, msg interface{}) {
	ctx.FireRead(msg)
}

func (ih *inouthandler) OnEvent(ctx *core.ChannelContext, event interface{}) {
	log.Infof("OnEvent: %+v", event)
	ctx.FireEvent(event)
}

// OnConnect called when new channel connected.
func (ih *inouthandler) OnConnect(ctx *core.ChannelContext, channel core.Channel) {
	ctx.FireConnect(channel)
}

// OnDisconnect called when channel disconnected.
func (ih *inouthandler) OnDisconnect(ctx *core.ChannelContext) {
	log.Infof("server disconnected: %+v", ctx)
	ctx.FireDisconnect()
}

func (ih *inouthandler) OnError(ctx *core.ChannelContext, err error) {
	log.Infof("OnError: %+v", err)
	ctx.FireError(err)
}

func TestPipeline(t *testing.T) {
	cp := core.NewChannelPipeline(nil)
	cp.AddLast(nil, "inhandler1", &inhandler1{})
	cp.AddLast(nil, "inhandler2", &inhandler2{})
	cp.AddFirst(nil, "inhandler1", &inhandler1{})

	cp.AddLast(nil, "outhandler", &outhandler{})
	cp.AddLast(nil, "inouthandler", &inouthandler{})
	cp.AddFirst(nil, "inouthandler", &inouthandler{})

	// cp.InvokeRead(&Msg{Hello: "world"})
	cp.FireWrite(&Msg{Hello: "world"})
}
