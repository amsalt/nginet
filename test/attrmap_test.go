package test

import (
	"testing"

	"github.com/amsalt/nginet/core"
	"github.com/amsalt/log"
)

func TestAttrMap(t *testing.T) {
	attr := core.NewAttrMap()
	for i := 0; i < 100; i++ {
		go attr.SetValue("test", i)
	}

	log.Infof("attr value: %+v\n", attr.IntValue("test"))
}
