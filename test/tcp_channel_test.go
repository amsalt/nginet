package test

import (
	"net"
	"testing"
	"time"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/aio"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/core/tcp"
	"github.com/amsalt/nginet/encoding"
	"github.com/amsalt/nginet/encoding/json"
	"github.com/amsalt/nginet/handler"
	"github.com/amsalt/nginet/message"
	"github.com/amsalt/nginet/message/idparser"
	"github.com/amsalt/nginet/message/packet"
)

type tcpChannel struct {
	Msg string
}

type tcpChannel2 struct {
	Test string
}

func TestTCPChannel(t *testing.T) {
	evtloop := aio.NewEventLoop()
	evtloop.Start()

	register := message.NewRegister()
	register.RegisterMsgByID(1, &tcpChannel{})
	register.RegisterMsgByID(2, &tcpChannel2{})

	parser := idparser.NewUint16ID()
	codec := encoding.MustGetCodec(json.CodecJSON)

	packetIdParser := handler.NewIDParser(register, parser)
	messageSerializer := handler.NewMessageSerializer(register, codec)
	messageDeserializer := handler.NewMessageDeserializer(register, codec)

	processMgr := message.NewProcessorMgr(register)
	processMgr.RegisterProcessor(&tcpChannel{}, func(ctx *core.ChannelContext, msg interface{}, args ...interface{}) {
		if m, ok := msg.([]byte); ok {
			log.Infof("tcpChannel handler: %+v", string(m))
		} else {
			log.Infof("tcpChannel handler: %+v", msg)
		}

		byteArr, err := codec.Marshal(&tcpChannel{Msg: "tcpChannel handler response data"})
		if err == nil {
			ctx.Write(packet.NewRawPacket(1, byteArr))
			ctx.Close()
		}
	})

	s := core.GetAcceptorBuilder(core.TCPServBuilder).Build(
		tcp.WithReadBufSize(1024),
		tcp.WithWriteBufSize(1024),
		tcp.WithMaxConnNum(100),
	)

	s.Pipeline().AddLast(nil, "inhandler1", &inhandler1{})
	s.InitSubChannel(func(channel core.SubChannel) {
		log.Infof("new channel created, channelId is %+v", channel.ID())
		channel.Pipeline().AddLast(nil, "IdleStateHandler", handler.NewIdleStateHandler(5, 5, false))
		channel.Pipeline().AddLast(nil, "inhandler1", &inhandler1{})
		channel.Pipeline().AddLast(nil, "PacketLengthDecoder", handler.NewPacketLengthDecoder(2))
		channel.Pipeline().AddLast(nil, "PacketLengthPrepender", handler.NewPacketLengthPrepender(2))
		channel.Pipeline().AddLast(nil, "rc4", handler.NewRc4Cipher("example"))
		channel.Pipeline().AddLast(nil, "MessageEncoder", handler.NewMessageEncoder(messageSerializer, packetIdParser))
		channel.Pipeline().AddLast(nil, "MessageDecoder", handler.NewMessageDecoder(messageDeserializer, packetIdParser))
		channel.Pipeline().AddLast(evtloop, "processor", handler.NewDefaultMessageHandler(processMgr))
	})
	addr, err := net.ResolveTCPAddr("tcp", ":7878")
	if err != nil {
		panic("bad net addr")
	}

	s.Listen(addr)
	go s.Accept()
	log.Infof("server local address:%+v", s.LocalAddr().String())

	c := tcp.NewClientChannel(&tcp.Options{WriteBufSize: 1024})
	c.InitSubChannel(func(channel core.SubChannel) {
		channel.Pipeline().AddLast(nil, "PacketLengthDecoder", handler.NewPacketLengthDecoder(2))
		channel.Pipeline().AddLast(nil, "PacketLengthPrepender", handler.NewPacketLengthPrepender(2))
		channel.Pipeline().AddLast(nil, "rc4", handler.NewRc4Cipher("example"))
		channel.Pipeline().AddLast(nil, "MessageEncoder", handler.NewMessageEncoder(messageSerializer, packetIdParser))
		channel.Pipeline().AddLast(nil, "MessageDecoder", handler.NewMessageDecoder(messageDeserializer, packetIdParser))
		channel.Pipeline().AddLast(nil, "processor", handler.NewDefaultMessageHandler(processMgr))
	})
	c.Connect(addr)
	log.Infof("client remote address:%+v", c.RemoteAddr().String())

	time.Sleep(2 * time.Second)
	c.Write(&tcpChannel{Msg: "tcp channel handler example"})
	time.Sleep(15 * time.Second)
	c.Close()
}
