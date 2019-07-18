package message

import (
	"encoding/binary"
	"reflect"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/encoding"
)

// package message defines all kinds of messages used in nginet.

// Packet represents the serialized form of a network message.
//
// the structure of a packet defined in nginet:
// 	- Length: Represents the length of the whole packet in byte array.
//  - PacketID: The application layer defines the mapping of packet ID, one packetID one network message.
//  - msg_payload: The payload of the packet.
//
// 					  Packet
// 	--------------|-------------|---------------
// 	|   Length    |   PacketID  |  msg_payload |
// 	--------------|-------------|---------------
//
type Packet interface {
	ID() interface{}
	Payload() interface{}
}

// Meta stores the information of a application layer's registered message.
// such as packetID, codec and the type information of stored message.
type Meta interface {
	// ID return message ID
	ID() interface{}

	// CreateInstance returns the zero value as an interface{} of the Meta's type
	CreateInstance() interface{}

	// SetCodec enable special codec for some protocol message.
	// If an codec is specified for meta,it will use the setted codec and ignore common codec setting
	SetCodec(c encoding.Codec)

	// Codec returns setted codec
	// It returns nil if codec not setted.
	Codec() encoding.Codec
}

// metaData is an implementation of Meta.
type metaData struct {
	msgID   interface{}
	msgName string
	msgType reflect.Type
	codec   encoding.Codec
}

func newMetaData(id interface{}, name string, msgType reflect.Type) *metaData {
	return &metaData{msgID: id, msgName: name, msgType: msgType}
}

// GetMsgID return the msg id of the meta data.
func (md *metaData) ID() interface{} {
	return md.msgID
}

// NewType return the interface of the meta data.
func (md *metaData) CreateInstance() interface{} {
	return reflect.New(md.msgType.Elem()).Interface()
}

func (md *metaData) SetCodec(codec encoding.Codec) {
	md.codec = codec
}

func (md *metaData) Codec() encoding.Codec {
	return md.codec
}

type Register interface {
	RegisterMsg(msg interface{}) (meta Meta)
	RegisterMsgByID(assignID interface{}, msg interface{}) Meta
	GetMetaByMsg(msg interface{}) Meta
	GetMetaByID(id interface{}) Meta
}

type ProcessorMgr interface {
	RegisterProcessor(msg interface{}, h ProcessorFunc) error
	RegisterProcessorByID(msgID interface{}, hf ProcessorFunc) error
	GetProcessorByID(msgID interface{}) *Processor
}

// PacketIDParser represents the parser to encode&decode a packet ID.
type PacketIDParser interface {
	// SetByteOrder specifies how to convert byte sequences into
	// 16-bit, 32-bit, or 64-bit unsigned integers.
	SetByteOrder(order binary.ByteOrder) PacketIDParser

	IDLen() int

	// Encode encodes the packet ID to bytes array.
	Encode(packetID interface{}, buf bytes.WriteOnlyBuffer) error

	// Decode decodes a packet ID from bytes array.
	Decode(rob bytes.ReadOnlyBuffer) (interface{}, error)
}
