package ws

import (
	"time"

	"github.com/amsalt/nginet/core"
)

func init() {
	core.Register(&wsServBuilder{})
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

func WithMaxHeaderSize(s int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).maxHeaderSize = s
	}
}

func WithTimeout(t time.Duration) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).timeout = t
	}
}

func WithCertFile(cf string) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).certFile = cf
	}
}

func WithKeyFile(kf string) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).keyFile = kf
	}
}

// WithAutoReconnect sets whether auto reconnect when disconnect with server.
func WithAutoReconnect(r bool) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).AutoReconnect = r
	}
}

// WithMaxReconnectTimes sets max number for trying reconnect with server.
func WithMaxReconnectTimes(n int) core.BuildOption {
	return func(o interface{}) {
		o.(*serverOptions).MaxReconnectTimes = n
	}
}

type wsServBuilder struct {
}

func (tb *wsServBuilder) Name() string {
	return core.WebsocketServBuilder
}

func (tb *wsServBuilder) Build(opt ...core.BuildOption) core.AcceptorChannel {
	opts := defaultServeroptions
	for _, o := range opt {
		o(&opts)
	}

	return newServerChannel(&opts)
}

var defaultCliOptions = Options{
	WriteBufSize:      1024,
	ReadBufSize:       1024,
	AutoReconnect:     false,
	MaxReconnectTimes: 20,
}

var defaultServeroptions = serverOptions{
	Options: &defaultCliOptions,

	maxConnNum: 1000 * 10000,
}

type Options struct {
	WriteBufSize      int
	ReadBufSize       int
	AutoReconnect     bool
	MaxReconnectTimes int
}

type serverOptions struct {
	*Options

	maxConnNum int

	maxHeaderSize int
	timeout       time.Duration

	// wss
	certFile string
	keyFile  string
}
