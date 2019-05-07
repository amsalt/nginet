package ws

import (
	"fmt"
	"net"
	"net/http"

	"github.com/amsalt/nginet/core"
	"github.com/amsalt/log"
	"github.com/gorilla/websocket"
)

type server struct {
	*core.AttrMap
	*core.Acceptor

	opts *serverOptions

	ln        net.Listener
	localAddr net.Addr
}

func newServerChannel(opts *serverOptions) core.AcceptorChannel {
	s := new(server)
	s.Acceptor = core.NewAcceptor()
	s.opts = opts

	return s
}

// Write writes message to opposite side.
func (s *server) Write(msg interface{}) error {
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
		log.Errorf("Websocket server init error: %+v", err)
	}
	s.ln = ln
}

// Accept accepts the next incoming call
func (s *server) Accept() {
	s.serve()
}

func (s *server) serve() {
	handler := http.NewServeMux()
	upgrader := websocket.Upgrader{
		HandshakeTimeout: s.opts.timeout,
		CheckOrigin:      func(_ *http.Request) bool { return true },
	}

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		if s.validate() {
			s.processNewConn(conn)
		}
	})

	httpServer := &http.Server{
		Handler:        handler,
		Addr:           s.localAddr.String(),
		MaxHeaderBytes: s.opts.maxHeaderSize,
		ReadTimeout:    s.opts.timeout,
		WriteTimeout:   s.opts.timeout,
	}

	if s.opts.certFile != "" || s.opts.keyFile != "" {
		go httpServer.ServeTLS(s.ln, s.opts.certFile, s.opts.keyFile)
	} else {
		go httpServer.Serve(s.ln)
	}
}

func (s *server) validate() bool {
	if len(s.SubChannels()) >= s.opts.maxConnNum {
		fmt.Errorf("too many connections, new connection will be Refused")
		return false
	}

	return true
}

func (s *server) processNewConn(conn *websocket.Conn) {
	log.Debugf("new connection: %+v", conn)
	s.FireConnect(core.NewDefaultSubChannel(newRawConn(conn), s.opts.ReadBufSize, s.opts.WriteBufSize))
}
