package proto

import (
	"errors"

	"github.com/amsalt/nginet/encoding"
	protobuf "github.com/gogo/protobuf/proto"
)

var (
	errInvalidMessageType = errors.New("proto: invalid message type, need protobuf.Message")
)

const (
	// CodecProtobuf represents a protobuf type.
	CodecProtobuf encoding.CodecType = "protobuf"
)

// Package proto defines the protobuf codec. Importing this package will
// register the codec.

func init() {
	encoding.RegisterCodec(newCodec())
}

// codec is a Codec implementation with protobuf.
type codec struct{}

func newCodec() *codec {
	return new(codec)
}

func (c *codec) Marshal(v interface{}) ([]byte, error) {
	pbmessage, ok := v.(protobuf.Message)
	if !ok {
		return nil, errInvalidMessageType
	}
	data, err := protobuf.Marshal(pbmessage)
	return data, err
}

func (c *codec) Unmarshal(data []byte, v interface{}) error {
	pbmessage, ok := v.(protobuf.Message)
	if !ok {
		return errInvalidMessageType
	}
	return protobuf.Unmarshal(data, pbmessage)
}

func (c *codec) Name() encoding.CodecType {
	return CodecProtobuf
}
