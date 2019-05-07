package test

import (
	"net"
	"testing"
	"time"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/aio"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/core/ws"
	"github.com/amsalt/nginet/encoding"
	"github.com/amsalt/nginet/encoding/json"
	"github.com/amsalt/nginet/handler"
	"github.com/amsalt/nginet/message"
	"github.com/amsalt/nginet/message/idparser"
	"github.com/amsalt/nginet/message/packet"
)

type wsChannel struct {
	Msg string
}

func TestWsChannel(t *testing.T) {
	evtloop := aio.NewEventLoop()
	evtloop.Start()

	register := message.NewRegister()
	register.RegisterMsgByID(1, &wsChannel{})

	parser := idparser.NewUint16ID()
	codec := encoding.MustGetCodec(json.CodecJSON)

	packetIdParser := handler.NewIDParser(register, parser)
	messageSerializer := handler.NewMessageSerializer(register, codec)
	messageDeserializer := handler.NewMessageDeserializer(register, codec)

	processMgr := message.NewProcessorMgr(register)
	processMgr.RegisterProcessorByID(1, func(ctx *core.ChannelContext, msg interface{}, args ...interface{}) {
		log.Infof("wsChannel handler: %+v", msg)
		byteArr, err := codec.Marshal(&wsChannel{Msg: "wsChannel handler response data"})
		if err == nil {
			ctx.Write(packet.NewRawPacket(1, byteArr))
		}
	})

	s := core.GetAcceptorBuilder(core.WebsocketServBuilder).Build(
		ws.WithWriteBufSize(1024),
	)
	s.Pipeline().AddLast(nil, "inhandler1", &inhandler1{})
	s.InitSubChannel(func(channel core.SubChannel) {
		log.Infof("new channel created, channelId is %+v", channel.ID())
		channel.Pipeline().AddLast(nil, "PacketLengthDecoder", handler.NewPacketLengthDecoder(2))
		channel.Pipeline().AddLast(nil, "PacketLengthPrepender", handler.NewPacketLengthPrepender(2))
		channel.Pipeline().AddLast(nil, "MessageSerializer", messageSerializer)
		channel.Pipeline().AddLast(nil, "IDParser", packetIdParser)
		channel.Pipeline().AddLast(nil, "MessageDeserializer", messageDeserializer)
		channel.Pipeline().AddLast(evtloop, "processor", handler.NewDefaultMessageHandler(processMgr))
	})
	addr, err := net.ResolveTCPAddr("tcp", ":7879")
	if err != nil {
		panic("bad net addr")
	}

	s.Listen(addr)
	go s.Accept()
	log.Infof("server local address:%+v", s.LocalAddr().String())

	c := ws.NewClientChannel()
	c.InitSubChannel(func(channel core.SubChannel) {
		channel.Pipeline().AddLast(nil, "PacketLengthDecoder", handler.NewPacketLengthDecoder(2))
		channel.Pipeline().AddLast(nil, "PacketLengthPrepender", handler.NewPacketLengthPrepender(2))
		channel.Pipeline().AddLast(nil, "MessageEncoder", handler.NewMessageEncoder(messageSerializer, packetIdParser))
		channel.Pipeline().AddLast(nil, "MessageDecoder", handler.NewMessageDecoder(messageDeserializer, packetIdParser))
		channel.Pipeline().AddLast(nil, "processor", handler.NewDefaultMessageHandler(processMgr))
	})
	c.Connect("ws://localhost:7879/")
	log.Infof("client remote address:%+v", c.RemoteAddr().String())

	time.Sleep(2 * time.Second)
	c.Write(&wsChannel{Msg: "ws channel handler test"})
	time.Sleep(15 * time.Second)
	c.Close()
}
