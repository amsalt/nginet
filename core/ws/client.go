package ws

import (
	"net/http"

	"github.com/amsalt/nginet/core"
	"github.com/amsalt/log"
	"github.com/gorilla/websocket"
)

type client struct {
	*core.Connector

	conn     *websocket.Conn
	response *http.Response

	opts *Options
}

func NewClientChannel(opts ...*Options) core.ConnectorChannel {
	c := &client{}
	c.Connector = core.NewConnector()
	if len(opts) > 0 {
		c.opts = opts[0]
	} else {
		c.opts = &defaultCliOptions
	}

	return c
}

func (c *client) Connect(addr interface{}) {
	log.Debugf("ws Connect addr: %+v", addr)
	d := &websocket.Dialer{}

	// TODO: addr
	conn, response, err := d.Dial(addr.(string), nil)
	if err != nil {
		panic(err)
	}

	c.conn = conn
	c.response = response
	c.FireConnect(core.NewDefaultSubChannel(newRawConn(conn), c.opts.WriteBufSize, c.opts.ReadBufSize))
}
