package encoding

import "fmt"

// Package codec defines the interface for the codec, and
// functions to register and retrieve codecs.

// A CodecType is a string codec type.
type CodecType string

// Codec defines the interface that nginet uses to encode and decode messages. a Codec's
// methods can be called from concurrent goroutines.
type Codec interface {
	// Marshal returns the special codec format of v.
	Marshal(v interface{}) ([]byte, error)

	// Unmarshal parses the special codec format into v.
	Unmarshal(data []byte, v interface{}) error

	// Name return the codec's name.
	Name() CodecType
}

var registeredCodecs map[CodecType]Codec

func init() {
	registeredCodecs = make(map[CodecType]Codec)
}

// RegisterCodec registers new Codec processor.
func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register nil codec")
	}

	codecName := codec.Name()
	if codecName == "" {
		panic("cannot register codec with empty name")
	}
	registeredCodecs[codecName] = codec
}

// GetCodec return registered codec by name.
// Note: if codec with name 'name' not registered,
// a nil codec will returned.
func GetCodec(name CodecType) Codec {
	return registeredCodecs[name]
}

// MustGetCodec the same with GetCodec.
// but it will panic when the codec of name 'name' not registred.
func MustGetCodec(name CodecType) Codec {
	c, exist := registeredCodecs[name]
	if !exist {
		panic(fmt.Sprintf("codec with name %v not registered", name))
	}
	return c
}
