package test

import (
	"testing"

	"github.com/amsalt/nginet/pool"
	"github.com/amsalt/log"
)

func BenchmarkPool(b *testing.B) {
	pool.SetSize(10)

	for i := 0; i < b.N; i++ {
		pool.Execute(func() {
			log.Infof("pool execte.")
		})
	}

}
