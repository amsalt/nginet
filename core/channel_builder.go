package core

import "strings"

var (
	// m is a map from name to channel builder.
	m = make(map[string]Builder)
)

const (
	// TCPServBuilder represents tcp server builder
	TCPServBuilder = "s_tcp"

	// WebsocketServBuilder represents websocket server builder
	WebsocketServBuilder = "s_websocket"

	// UDPServBuilder represents udp server builder
	UDPServBuilder = "s_udp"
)

// BuildOption represents transport builder option.
type BuildOption func(interface{})

// Register registers the transport(e.g. TcpAcceptor) builder to the transport map. b.Name
// (lowercased) will be used as the name registered with this builder.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe. If multiple Acceptor are
// registered with the same name, the one registered last will take effect.
func Register(b Builder) {
	m[strings.ToLower(b.Name())] = b
}

// Builder builds transport.
type Builder interface {
	Build(opts ...BuildOption) AcceptorChannel
	Name() string
}

// GetAcceptorBuilder returns acceptor builder by name, ignore case.
func GetAcceptorBuilder(n string) Builder {
	return m[strings.ToLower(n)]
}
