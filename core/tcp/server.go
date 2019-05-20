package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/core"
)

// server represents a tcp server.
// implements the AcceptorChannel interface.
type server struct {
	opts *serverOptions
	*core.AttrMap
	*core.Acceptor

	ln        net.Listener
	localAddr net.Addr

	retryDelay time.Duration
}

func newServerChannel(opts *serverOptions) core.AcceptorChannel {
	s := new(server)
	s.opts = opts
	s.Acceptor = core.NewAcceptor()

	return s
}

// Write writes message to opposite side.
func (s *server) Write(msg interface{}, extra ...interface{}) error {
	// nothing to do.
	return nil
}

// LocalAddr returns the local addr.
func (s *server) LocalAddr() net.Addr {
	return s.localAddr
}

// RemoteAddr return the opposite side addr.
func (s *server) RemoteAddr() net.Addr {
	panic("not implementation")
}

// Listen announces on the local network address.
func (s *server) Listen(addr net.Addr) {
	s.localAddr = addr

	ln, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		panic(fmt.Errorf("TCP server init error: %+v", err))
	}
	s.ln = ln
}

// Accept accepts the next incoming call
func (s *server) Accept() {
	s.accept()
}

func (s *server) accept() {
	for {
		conn, err := s.ln.(*net.TCPListener).AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				s.waitWhenTemporaryErr()
				continue
			}
			// Stop
			return
		}
		s.resetWhenSucc()

		if !s.validate(conn) {
			conn.Close()
			continue
		}

		s.applyOptions(conn)
		s.processNewConn(conn)
	}
}

func (s *server) Close() {
	s.ln.Close()
}

func (s *server) waitWhenTemporaryErr() {
	if s.retryDelay == 0 {
		s.retryDelay = 5 * time.Millisecond
	} else {
		s.retryDelay *= 2
	}
	if max := 1 * time.Second; s.retryDelay > max {
		s.retryDelay = max
	}

	time.Sleep(s.retryDelay)
}

func (s *server) resetWhenSucc() {
	s.retryDelay = 0
}

func (s *server) validate(conn *net.TCPConn) bool {
	if len(s.SubChannels()) > s.opts.maxConnNum {
		return false
	}
	return true
}

func (s *server) applyOptions(conn *net.TCPConn) {
	conn.SetNoDelay(s.opts.noDelay)
	conn.SetKeepAlive(s.opts.keepalive)

	if s.opts.tcpReadBufSize > 0 {
		conn.SetReadBuffer(s.opts.tcpReadBufSize)
	}

	if s.opts.tcpWriteBufSize > 0 {
		conn.SetWriteBuffer(s.opts.tcpWriteBufSize)
	}

	if s.opts.keepalivePeriod > 0 {
		conn.SetKeepAlivePeriod(time.Second * time.Duration(s.opts.keepalivePeriod))
	}

	if s.opts.linger > 0 {
		conn.SetLinger(s.opts.linger)
	}
}

func (s *server) processNewConn(conn *net.TCPConn) {
	log.Debugf("new connection: %+v", conn)
	s.FireConnect(core.NewDefaultSubChannel(newRawConn(conn), s.opts.ReadBufSize, s.opts.WriteBufSize))
}
