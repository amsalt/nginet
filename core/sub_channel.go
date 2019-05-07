package core

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/amsalt/nginet/bytes"

	"github.com/amsalt/log"
)

var (
	ErrWriteMsgQueueFull = errors.New("write queue is full")
	ErrConnLost          = errors.New("connection lost")
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
}

// NewDefaultSubChannel returns a new instance of SubChannel
// The same time, the subchannel will start the read loop and write loop
// to serve reading message and writting message.
func NewDefaultSubChannel(conn RawConn, readBufSize, writeBufSize int) SubChannel {
	dsc := &DefaultSubChannel{conn: conn}
	dsc.BaseChannel = NewBaseChannel(dsc)
	dsc.closeChan = make(chan byte)

	dsc.writeBuf = make(chan interface{}, writeBufSize)
	dsc.readBufSize = readBufSize

	dsc.start()
	return dsc
}

func (dsc *DefaultSubChannel) start() {
	go dsc.readloop()
	go dsc.writeloop()
}

func (dsc *DefaultSubChannel) readloop() {
	log.Debug("start read loop.")
	readerBuf := bytes.NewReadOnlyBuffer(dsc.readBufSize)

	timer := time.NewTimer(5 * time.Second)

	for {
		err := dsc.conn.Read(readerBuf)
		if err != nil {
			goto ERR
		}

		// ensure process all messages.
		for readerBuf.Len() > 0 {
			oldBufLen := readerBuf.Len()
			dsc.FireRead(readerBuf)
			if oldBufLen == readerBuf.Len() {
				break
			}
		}

		// validate
		select {
		case <-timer.C:
			if !dsc.isAlive() {
				dsc.Close()
			}
			timer.Reset(5 * time.Second)
		case <-dsc.closeChan:
			timer.Stop()
			goto CLOSED
		default:
		}
	}
ERR:
	dsc.Close()
	return
CLOSED:
	// exit the read loop.
}

func (dsc *DefaultSubChannel) isAlive() bool {
	return true
}

func (dsc *DefaultSubChannel) writeloop() {
	log.Debug("start write loop.")
	for msg := range dsc.writeBuf {
		if msg == nil {
			continue
		}
		dsc.Pipeline().FireWrite(msg)
	}
	dsc.Close()
}

// Write writes message to opposite side.
func (dsc *DefaultSubChannel) Write(msg interface{}) (err error) {
	select {
	case dsc.writeBuf <- msg:
	case <-dsc.closeChan:
		err = ErrConnLost
	default:
		err = ErrWriteMsgQueueFull
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

// RawConn returns the raw connection.
func (dsc *DefaultSubChannel) RawConn() RawConn {
	return dsc.conn
}
