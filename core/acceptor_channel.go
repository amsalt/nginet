package core

import (
	"sync"

	"github.com/amsalt/log"
)

// Acceptor represents a simple base server for share code.
type Acceptor struct {
	sync.RWMutex
	*BaseChannel

	// handle new subchannel.
	initCb InitChannelCb

	subChannels []SubChannel
	channels    map[interface{}]SubChannel
}

// NewAcceptor create a Acceptor instance which can accept new connection from client.
func NewAcceptor() *Acceptor {
	acceptor := new(Acceptor)
	acceptor.BaseChannel = NewBaseChannel(acceptor)
	acceptor.channels = make(map[interface{}]SubChannel)
	return acceptor
}

// InitSubChannel sets the callback when create a new SubChannel to init it.
func (acceptor *Acceptor) InitSubChannel(sub InitChannelCb) {
	acceptor.initCb = sub
}

func (acceptor *Acceptor) SubChannelInitializer() InitChannelCb {
	return acceptor.initCb
}

func (acceptor *Acceptor) initChannel(c SubChannel) {
	acceptor.initCb(c)
}

// FireConnect fires a Connect event.
func (acceptor *Acceptor) FireConnect(channel Channel) InboundInvoker {
	subChannel, ok := channel.(SubChannel)
	if !ok {
		log.Errorf("Acceptor.FireConnect channel should be SubChannel")
		return nil
	}
	acceptor.initChannel(subChannel)

	acceptor.Lock()
	acceptor.subChannels = append(acceptor.subChannels, subChannel)
	acceptor.channels[subChannel.ID()] = subChannel
	acceptor.Unlock()
	return subChannel.Pipeline().FireConnect(subChannel)
}

// FireDisconnect fires a Disconnect event.
func (acceptor *Acceptor) FireDisconnect() InboundInvoker {
	return acceptor.Pipeline().FireDisconnect()
}

// FireRead fires a read event.
func (acceptor *Acceptor) FireRead(msg interface{}) InboundInvoker {
	panic("Acceptor.FireRead not implement.")
}

// FireOnEvent fires a  event.
func (acceptor *Acceptor) FireOnEvent(event interface{}) InboundInvoker {
	return acceptor.Pipeline().FireOnEvent(event)
}

func (acceptor *Acceptor) FireError(err error) InboundInvoker {
	return acceptor.Pipeline().FireError(err)
}

// ------------------- Channel Manager methods -------------------

// SubChannels returns all SubChannels belong to this AcceptorChannel.
func (acceptor *Acceptor) SubChannels() []SubChannel {
	return acceptor.subChannels
}

// Broadcast broadcasts message to all client channels.
// Just a simple wrapper for connection manager.
// For better performance, should only do FilterPipeline once for all Channel
// and call RawConn().Write(msg) for less filter operations and less memory cost.
func (acceptor *Acceptor) Broadcast(msg interface{}) error {
	log.Debugf("acceptor subchannels: %+v", acceptor.subChannels)
	for _, channel := range acceptor.subChannels {
		err := channel.Write(msg)
		if err != nil {
			log.Errorf("Acceptor.Broadcast failed: %+v", err)
		}
	}
	return nil
}

// Multicast sends message to specified channels.
func (acceptor *Acceptor) Multicast(msg interface{}, channelIDs []interface{}) error {
	acceptor.RLock()
	defer acceptor.RUnlock()
	for _, id := range channelIDs {
		acceptor.channels[id].Write(msg)
	}

	return nil
}
