package tcp

import (
	"github.com/amsalt/nginet/core"
)

func init() {
	core.Register(&tcpServBuilder{})
}

// WithWriteBufSize sets max size of pending wirte.
func WithWriteBufSize(s int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).WriteBufSize = s
	}
}

// WithReadBufSize sets max size of pending read.
func WithReadBufSize(s int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).ReadBufSize = s
	}
}

// WithMaxConnNum sets max connection number.
func WithMaxConnNum(mcn int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).maxConnNum = mcn
	}
}

// WithTCPWriteBufSize sets the size of the operating system's
// transmit buffer associated with the connection.
func WithTCPWriteBufSize(s int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).tcpWriteBufSize = s
	}
}

// WithTCPReadBufSize sets the size of the operating system's
// receive buffer associated with the connection.
func WithTCPReadBufSize(s int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).tcpReadBufSize = s
	}
}

// WithNodelay controls whether the operating system should delay
// packet transmission in hopes of sending fewer packets (Nagle's
// algorithm).  The default is true (no delay), meaning that data is
// sent as soon as possible after a Write.
func WithNodelay(b bool) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).noDelay = b
	}
}

// WithKeepalive sets whether the operating system should send
// keepalive messages on the connection.
func WithKeepalive(b bool) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).keepalive = b
	}
}

// WithKeepalivePeriod sets period between keep alives.
func WithKeepalivePeriod(sec int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).keepalivePeriod = sec
	}
}

// WithLinger sets the behavior of Close on a connection which still
// has data waiting to be sent or to be acknowledged.
//
// If sec < 0 (the default), the operating system finishes sending the
// data in the background.
//
// If sec == 0, the operating system discards any unsent or
// unacknowledged data.
//
// If sec > 0, the data is sent in the background as with sec < 0. On
// some operating systems after sec seconds have elapsed any remaining
// unsent data may be discarded.
func WithLinger(sec int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).linger = sec
	}
}

type tcpServBuilder struct {
}

func (tb *tcpServBuilder) Name() string {
	return core.TCPServBuilder
}

func (tb *tcpServBuilder) Build(opt ...core.BuildOption) core.AcceptorChannel {
	opts := defaultServeroptions
	for _, o := range opt {
		o(&opts)
	}

	return newServerChannel(&opts)
}

var defaultCliOptions = Options{
	WriteBufSize:  1024,
	ReadBufSize:   1024,
	AutoReconnect: false,
}

var defaultServeroptions = serverOptions{
	Options: &defaultCliOptions,

	maxConnNum: 1000 * 10000,
	noDelay:    true,
	keepalive:  false,
}

type Options struct {
	WriteBufSize  int
	ReadBufSize   int
	AutoReconnect bool
}

type serverOptions struct {
	*Options

	maxConnNum      int
	tcpWriteBufSize int
	tcpReadBufSize  int

	noDelay         bool
	keepalive       bool
	keepalivePeriod int
	linger          int
}
