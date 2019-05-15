package test

import (
	"net"
	"testing"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/core/tcp"
	"github.com/amsalt/nginet/gnetlog"
	"github.com/amsalt/nginet/handler"
)

func init() {
	gnetlog.Init()
}

type texthandler struct {
}

func (ih *texthandler) OnRead(ctx *core.ChannelContext, msg interface{}) {
	log.Infof("texthandler msg: %+v", msg)
	ctx.Write([]byte("response message.\n"))
	ctx.FireRead(msg)
}

// OnConnect called when new channel connected.
func (ih *texthandler) OnConnect(ctx *core.ChannelContext, channel core.Channel) {
	log.Infof("new client connected: %+v", ctx)
	ctx.FireConnect(channel)
}

func (ih *texthandler) OnError(ctx *core.ChannelContext, err error) {
	log.Infof("OnError: %+v", err)
	ctx.FireError(err)
}

func (ih *texthandler) OnEvent(ctx *core.ChannelContext, event interface{}) {
	log.Infof("OnEvent: %+v", event)
	ctx.FireEvent(event)
	ctx.Close()
}

// OnDisconnect called when channel disconnected.
func (ih *texthandler) OnDisconnect(ctx *core.ChannelContext) {
	log.Infof("client disconnected: %+v", ctx)
	ctx.FireDisconnect()
}
func TestTelnet(t *testing.T) {
	s := core.GetAcceptorBuilder(core.TCPServBuilder).Build(
		tcp.WithReadBufSize(1024),
		tcp.WithWriteBufSize(1024),
		tcp.WithMaxConnNum(100),
	)

	s.InitSubChannel(func(channel core.SubChannel) {
		channel.Pipeline().AddLast(nil, "stringencoder", handler.NewStringEncoder())
		channel.Pipeline().AddLast(nil, "texthandler", &texthandler{})
	})
	addr, err := net.ResolveTCPAddr("tcp", ":7878")
	if err != nil {
		panic("bad net addr")
	}

	s.Listen(addr)
	s.Accept()
}
