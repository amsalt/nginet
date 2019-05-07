package core

import (
	"net"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/shortid"
)

// Channel is a nexus to a network socket or a component which is capable of I/O
// operations such as read, write, connect, and close.
type Channel interface {
	ID() interface{}

	// Attr returns the AttrMap.
	// each channel contains a AttrMap
	Attr() *AttrMap

	// Write writes message to opposite side.
	Write(msg interface{}) error

	// Close close the connection
	Close()

	// Pipeline returns the ChannelPipeline.
	Pipeline() ChannelPipeline

	// LocalAddr returns the local addr.
	LocalAddr() net.Addr

	// RemoteAddr return the opposite side addr.
	RemoteAddr() net.Addr

	RawConn() RawConn
}

// InitChannelCb defines the behaviors when AcceptorChannel creates a new sub channel.
type InitChannelCb func(channel SubChannel)

// AcceptorChannel represents a server-side socket.
type AcceptorChannel interface {
	Channel
	InboundInvoker
	ChannelMgr

	// InitSubChannel set the InitChannel function which defines the behaviors when new sub channel created.
	InitSubChannel(sub InitChannelCb)

	// Listen announces on the local network address.
	Listen(addr net.Addr)

	// Accept accepts the next incoming call
	Accept()
}

// ChannelMgr represents a manager of Channel.
type ChannelMgr interface {
	// SubChannels returns all SubChannels belong to this AcceptorChannel.
	SubChannels() []SubChannel

	// Broadcast broadcasts message to all client channels.
	Broadcast(msg interface{}) error

	// Multicast sends message to specified channels.
	Multicast(msg interface{}, channelIDs []interface{}) error
}

// ConnectorChannel respresents a client-side socket.
type ConnectorChannel interface {
	Channel
	InboundInvoker

	// InitSubChannel set the InitChannel function which defines the behaviors when new sub channel created.
	InitSubChannel(sub InitChannelCb)

	// Connect connects to the special address.
	Connect(addr interface{})
}

// SubChannel represents a server-side connection connected by a client.
type SubChannel interface {
	Channel
	InboundInvoker

	// GracefullyClose closes gracefully with all message sent before close.
	GracefullyClose()
}

type defaultIDGenerator struct {
}

func newDefaultIDGenerator() ChannelIDGenerator {
	return new(defaultIDGenerator)
}

func (dg *defaultIDGenerator) GenID() interface{} {
	cid, err := shortid.Generate()
	if err != nil {
		log.Debug("CreateChannel shortid.Generate failed %+v", err)
		return nil
	}
	return cid
}

// BaseChannel represents a basic implementation of Channel.
// it is used as the super class of other Channel.
type BaseChannel struct {
	AttrMap
	pipeline  ChannelPipeline
	id        interface{}
	generator ChannelIDGenerator
}

// NewBaseChannel creates a BaseChannel instance.
func NewBaseChannel(sub Channel) *BaseChannel {
	bc := &BaseChannel{}
	bc.pipeline = NewChannelPipeline(sub)
	bc.id = bc.genID()
	return bc
}

func (bc *BaseChannel) genID() interface{} {
	if bc.generator == nil {
		return newDefaultIDGenerator().GenID()
	}
	return bc.generator.GenID()
}

func (bc *BaseChannel) ID() interface{} {
	return bc.id
}

// Attr returns the AttrMap which contains a sets of attrs.
func (bc *BaseChannel) Attr() *AttrMap {
	return &bc.AttrMap
}

// Pipeline returns the ChannelPipeline.
func (bc *BaseChannel) Pipeline() ChannelPipeline {
	return bc.pipeline
}

// LocalAddr returns the local addr.
func (bc *BaseChannel) LocalAddr() net.Addr {
	panic("BaseChannel LocalAddr not implement, need be overrided")
}

// RemoteAddr return the opposite side addr.
func (bc *BaseChannel) RemoteAddr() net.Addr {
	panic("BaseChannel RemoteAddr not implement, need be overrided")
}

func (bc *BaseChannel) Write(msg interface{}) error {
	panic("BaseChannel Write not implement, need be overrided")
}

func (bc *BaseChannel) Close() {
	panic("BaseChannel Close not implement, need be overrided")
}

func (bc *BaseChannel) RawConn() RawConn {
	panic("BaseChannel RawConn not implement, need be overrided")
}
