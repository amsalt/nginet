package json

import (
	"encoding/json"

	"github.com/amsalt/nginet/encoding"
)

// Package json defines the JSON codec. Importing this package will
// register the codec.

const (
	// CodecJSON represents a json type.
	CodecJSON encoding.CodecType = "json"
)

func init() {
	encoding.RegisterCodec(newCodec())
}

// codec is a Codec implementation with json.
type codec struct{}

func newCodec() *codec {
	return new(codec)
}

func (c *codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (c *codec) Name() encoding.CodecType {
	return CodecJSON
}
