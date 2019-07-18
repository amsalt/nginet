package test

import (
	"testing"

	"github.com/amsalt/nginet/encoding"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/encoding/json"
	_ "github.com/amsalt/nginet/encoding/proto"
)

type msg struct {
	Hello string `json:"hello"`
}

func TestEncoding(t *testing.T) {
	codecJSON := encoding.GetCodec(json.CodecJSON)
	m := &msg{Hello: "world"}
	b, err := codecJSON.Marshal(m)
	if err != nil {
		log.Errorf("TestEncoding failed (%+v)", err)
	}
	log.Infof("TestEncoding result (%+v)", string(b))
}
