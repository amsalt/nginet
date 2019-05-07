package core

import (
	"net"

	"github.com/amsalt/log"
)

type Connector struct {
	*BaseChannel
	initCb     InitChannelCb
	subChannel SubChannel
}

func NewConnector() *Connector {
	c := new(Connector)
	c.BaseChannel = NewBaseChannel(c)
	return c
}

func (c *Connector) InitSubChannel(sub InitChannelCb) {
	c.initCb = sub
}

func (c *Connector) initChannel(channel SubChannel) {
	c.initCb(channel)
	log.Debugf("Connector.FireConnect init channel %+v", channel)
	c.subChannel = channel
}

func (c *Connector) FireConnect(channel Channel) InboundInvoker {
	subChannel, ok := channel.(SubChannel)
	if !ok {
		log.Errorf("Connector.FireConnect type err, channel is %T", channel)
		return nil
	}
	c.initChannel(subChannel)

	return channel.Pipeline().FireConnect(channel)
}
func (c *Connector) FireDisconnect() InboundInvoker {
	return c.Pipeline().FireDisconnect()
}

func (c *Connector) FireRead(msg interface{}) InboundInvoker {
	panic("not implement.")
}

func (c *Connector) FireError(err error) InboundInvoker {
	return c.Pipeline().FireError(err)
}

// Wrapper method for SubChannel

func (c *Connector) Write(msg interface{}) error {
	return c.subChannel.Write(msg)
}

// LocalAddr returns the local addr.
func (c *Connector) LocalAddr() net.Addr {
	return c.subChannel.LocalAddr()
}

// RemoteAddr return the opposite side addr.
func (c *Connector) RemoteAddr() net.Addr {
	return c.subChannel.RemoteAddr()
}

func (c *Connector) Close() {
	c.subChannel.Close()
}
