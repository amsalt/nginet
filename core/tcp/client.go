package tcp

import (
	"net"

	"github.com/amsalt/nginet/core"
)

type client struct {
	opts *Options
	*core.Connector
	conn net.Conn
	addr net.Addr
}

func NewClientChannel(opts ...*Options) core.ConnectorChannel {
	c := &client{}
	c.Connector = core.NewConnector()
	if len(opts) == 0 {
		c.opts = &defaultCliOptions
	} else {
		c.opts = opts[0]
	}

	return c
}

func (c *client) Connect(addr interface{}) (core.SubChannel, error) {
	netaddr, ok := addr.(net.Addr)
	if !ok {
		panic("tcp.client connect option must be net.Addr type")
	}
	c.addr = netaddr
	conn, err := net.Dial(netaddr.Network(), netaddr.String())
	if err != nil {
		return nil, err
	}

	subChannel := core.NewDefaultSubChannel(newRawConn(conn), c.opts.ReadBufSize, c.opts.WriteBufSize, &core.ReconnectOpts{AutoReconnect: c.opts.AutoReconnect, MaxReconnectTimes: c.opts.MaxReconnectTimes})
	c.FireConnect(subChannel)
	return subChannel, nil
}

func (c *client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
