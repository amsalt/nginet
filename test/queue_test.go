package test

import (
	"sync"
	"testing"

	"github.com/amsalt/nginet/aio"
)

const (
	Number = 10000000
)

func BenchmarkQueue(t *testing.B) {
	var wg sync.WaitGroup
	wg.Add(Number)
	queue := aio.NewQueue()
	for i := 0; i < Number; i++ {
		queue.Push(i)
	}

	for j := 0; j < 1000; j++ {
		go func() {
			var copy []interface{}
			queue.PollAll(&copy)
			wg.Add(-len(copy))
		}()
	}

	wg.Wait()

}

func BenchmarkChanQueue(t *testing.B) {
	var wg sync.WaitGroup
	wg.Add(Number)
	queue := aio.NewChanQueue(Number*2, -1)
	for i := 0; i < Number; i++ {
		queue.Push(i)
	}

	for j := 0; j < Number; j++ {
		go func() {
			queue.Get()
			wg.Add(-1)
		}()
	}

	wg.Wait()
}
