package core

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/amsalt/nginet/bytes"

	"github.com/amsalt/log"
)

const (
	// RetryMaxWaitSec represents the max wait time for retry reconnect.
	RetryMaxWaitSec = 120
)

var (
	// ErrWriteMsgQueueFull represents an error that the reserved queue is full.
	ErrWriteMsgQueueFull = errors.New("write queue is full")

	// ErrConnLost connection lost
	ErrConnLost = errors.New("connection lost")
)

// DefaultSubChannel represents a default implementation of SubChannel.
// When a new connection conntect to ServreChannel, a corresponding
//  SubChannel will be created.
type DefaultSubChannel struct {
	*BaseChannel
	sync.Mutex

	conn        RawConn
	closeChan   chan byte
	closeFlag   bool
	writeBuf    chan interface{}
	readBufSize int

	// if open, will reconnect with the same SubChannel instance
	autoReconnect     bool
	reconnectTimes    int
	maxReconnectTimes int
	reconnecting      bool
}

type ReconnectOpts struct {
	AutoReconnect     bool
	MaxReconnectTimes int
}

// NewDefaultSubChannel returns a new instance of SubChannel
// The same time, the subchannel will start the read loop and write loop
// to serve reading message and writting message.
func NewDefaultSubChannel(conn RawConn, readBufSize, writeBufSize int, reconnOpts ...*ReconnectOpts) SubChannel {
	dsc := &DefaultSubChannel{conn: conn}
	dsc.BaseChannel = NewBaseChannel(dsc)
	dsc.closeChan = make(chan byte)

	if len(reconnOpts) > 0 {
		dsc.autoReconnect = reconnOpts[0].AutoReconnect
		dsc.maxReconnectTimes = reconnOpts[0].MaxReconnectTimes
	}

	dsc.writeBuf = make(chan interface{}, writeBufSize)
	dsc.readBufSize = readBufSize

	dsc.start()
	return dsc
}

func (dsc *DefaultSubChannel) start() {
	go dsc.readloop()
	go dsc.writeloop()
}

func (dsc *DefaultSubChannel) resetReconn() {
	dsc.reconnecting = false
	dsc.reconnectTimes = 0
}

func (dsc *DefaultSubChannel) reconnect() bool {
	dsc.reconnecting = true
	for dsc.reconnectTimes < dsc.maxReconnectTimes {
		log.Warning("reconnecting...")
		dsc.reconnectTimes++
		netaddr := dsc.conn.RemoteAddr()

		conn, err := net.Dial(netaddr.Network(), netaddr.String())
		if err == nil {
			log.Debug("reconnect success: %+v", conn)
			dsc.conn.SetConn(conn)
			dsc.resetReconn()
			return true
		} else {
			time.Sleep((time.Duration(dsc.reconnectTimes) * 2 % RetryMaxWaitSec) * time.Second)
		}
	}

	dsc.resetReconn()
	return false
}

func (dsc *DefaultSubChannel) readloop() {
	log.Debug("start read loop.")
	readerBuf := bytes.NewReadOnlyBuffer(dsc.readBufSize)

	for {
		err := dsc.conn.Read(readerBuf)

		if err != nil {
			log.Infof("read message err: %+v", err)
			if !dsc.autoReconnect {
				goto ERR
			} else {
				reconnectSuccess := dsc.reconnect()
				if !reconnectSuccess {
					goto ERR
				}
			}
		}

		// ensure process all messages.
		for readerBuf.Len() > 0 {
			oldBufLen := readerBuf.Len()
			dsc.FireRead(readerBuf)
			if oldBufLen == readerBuf.Len() {
				break
			}
		}
	}
ERR:
	dsc.Close()
	return
}

func (dsc *DefaultSubChannel) writeloop() {
	log.Debug("start write loop.")
	for msg := range dsc.writeBuf {
		if msg == nil {
			continue
		}

		for dsc.reconnecting {
			log.Debug("wait write when reconnecting.")
			time.Sleep(time.Second)
		}

		dsc.Pipeline().FireWrite(msg)
	}
	dsc.Close()
}

// Write writes message to opposite side.
func (dsc *DefaultSubChannel) Write(msg interface{}, extra ...interface{}) (err error) {
	var output interface{}
	var combines []interface{}
	if len(extra) > 0 {
		combines = append(combines, msg, extra[0])
		output = combines
	} else {
		output = msg
	}

	select {
	case dsc.writeBuf <- output:
	case <-dsc.closeChan:
		err = ErrConnLost
	default:
		err = ErrWriteMsgQueueFull
	}

	if err != nil {
		log.Errorf("subchannel write message err: %+v", err)
	}
	return
}

// LocalAddr returns the local addr.
func (dsc *DefaultSubChannel) LocalAddr() net.Addr {
	if dsc.conn == nil {
		return nil
	}
	return dsc.conn.LocalAddr()
}

// RemoteAddr return the opposite side addr.
func (dsc *DefaultSubChannel) RemoteAddr() net.Addr {
	if dsc.conn == nil {
		return nil
	}
	return dsc.conn.RemoteAddr()
}

func (dsc *DefaultSubChannel) GracefullyClose() {

}

// Close closes the connection.
func (dsc *DefaultSubChannel) Close() {
	dsc.Lock()
	defer dsc.Unlock()

	if dsc.conn != nil && !dsc.closeFlag {
		dsc.closeFlag = true
		close(dsc.closeChan)
		dsc.conn.Close()
		dsc.FireDisconnect()
	}
}

// FireConnect fires a connect event. Nothing to do for a SubChannel.
func (dsc *DefaultSubChannel) FireConnect(channel Channel) InboundInvoker {
	// do nothing.
	return dsc
}

// FireDisconnect fires a disconnect event.
func (dsc *DefaultSubChannel) FireDisconnect() InboundInvoker {
	return dsc.Pipeline().FireDisconnect()
}

// FireRead fires a read event.
func (dsc *DefaultSubChannel) FireRead(msg interface{}) InboundInvoker {
	dsc.Pipeline().FireRead(msg)
	return dsc
}

func (dsc *DefaultSubChannel) FireError(err error) InboundInvoker {
	dsc.Pipeline().FireError(err)
	return dsc
}

func (dsc *DefaultSubChannel) FireEvent(event interface{}) InboundInvoker {
	dsc.Pipeline().FireEvent(event)
	return dsc
}

// RawConn returns the raw connection.
func (dsc *DefaultSubChannel) RawConn() RawConn {
	return dsc.conn
}
