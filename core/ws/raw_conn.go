package ws

import (
	"errors"
	"net"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
	"github.com/gorilla/websocket"
)

type rawConn struct {
	conn *websocket.Conn
}

func newRawConn(conn *websocket.Conn) core.RawConn {
	r := &rawConn{conn: conn}
	return r
}

// Write writes message to opposite side.
func (r *rawConn) Write(msg []byte) {
	if r.conn != nil {
		r.conn.WriteMessage(websocket.BinaryMessage, msg)
	}
}

func (r *rawConn) Read(buf bytes.ReadOnlyBuffer) error {
	if r.conn == nil {
		return errors.New("conn is nil")
	}

	_, reader, err := r.conn.NextReader()
	if err != nil {
		return err
	}
	_, err = buf.ReadFrom(reader)

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
	return errors.New("ws.rawConn Close() failed for conn is nil")
}
