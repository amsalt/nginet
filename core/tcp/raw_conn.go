package tcp

import (
	"errors"
	"net"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
)

type rawConn struct {
	conn net.Conn
}

func newRawConn(conn net.Conn) core.RawConn {
	r := &rawConn{conn: conn}
	return r
}

// Write writes message to opposite side.
func (r *rawConn) Write(msg []byte) {
	if r.conn != nil {
		r.conn.Write(msg)
	}
}

func (r *rawConn) Read(buf bytes.ReadOnlyBuffer) error {
	_, err := buf.ReadFrom(r.conn)
	return err
}

// LocalAddr returns the local addr.
func (r *rawConn) LocalAddr() net.Addr {
	if r.conn == nil {
		return nil
	}

	return r.conn.LocalAddr()
}

// RemoteAddr return the opposite side addr.
func (r *rawConn) RemoteAddr() net.Addr {
	if r.conn == nil {
		return nil
	}
	return r.conn.RemoteAddr()
}
func (r *rawConn) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return errors.New("tcp.rawConn Close() failed for conn is nil")
}
