package core

import (
	"net"

	"github.com/amsalt/nginet/bytes"
)

// package core implements the core functions of network.
// scheduler schedules network event(e.g. read new message) to filters and call handlers.
// handler represents a message logic processor.
//
// Channel represents a connector or a acceptor or a specific connection.
// ChannelContext is the context of the ChannelHandler of a specific Channel,
// 	which is used to provide interaction between Channel and ChannelHandler,
// 	and a ChannelContext for a specific Handler of each Channel.
// ChannelHandler is a message interceptor or processor. A channel may have N ChannelHandlers.
// ChannelPipeline is a message pipeline for concatenating ChannelHandlers.
// 	One ChannelPipeline per channel

// RawConn represents a raw connection, such as tcp connection or websocket connection.
// only contains basic API for read, write and close the raw connection.

type RawConn interface {
	Read(buf bytes.ReadOnlyBuffer) error
	Write(data []byte)

	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr
}

type ChannelIDGenerator interface {
	GenID() interface{}
}

type Executor interface {
	Execute(task func())
}

// InboundInvoker invokes inbound event handler.
type InboundInvoker interface {
	// FireConnect fire a connect event when new channel created.
	FireConnect(newChannel Channel) InboundInvoker

	// FireDisconnect fire a disconnect event when a channel destoryed.
	FireDisconnect() InboundInvoker

	// FireRead fire a read event when new data comes.
	FireRead(msg interface{}) InboundInvoker

	FireError(err error) InboundInvoker
}

// OutboundInvoker invokes outbound event handler.
type OutboundInvoker interface {
	// FireWrite fire a write event when send msg to channel.
	FireWrite(msg interface{}) OutboundInvoker
}
