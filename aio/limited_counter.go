package aio

import (
	"fmt"
	"sync"
)

var gLimitedCounter *LimitedCounter
var once sync.Once

const (
	defaultSize = 10240
)

// InitLimitedCounter set size of Goroutine counter.
func InitLimitedCounter(size int) {
	gLimitedCounter = NewLimitedCounter(size)
}

// Go calls task in new goroutine.
func Go(task func()) {
	if gLimitedCounter == nil {
		once.Do(func() {
			gLimitedCounter = NewLimitedCounter(defaultSize)
		})

	}

	gLimitedCounter.ApplyOne()
	go func() {
		defer gLimitedCounter.ReturnOne()
		task()
	}()
}

// LimitedCounter is a counter for limiting the max number of goroutines.
type LimitedCounter struct {
	counter chan int
	wg      sync.WaitGroup
}

func NewLimitedCounter(size int) *LimitedCounter {
	if size <= 0 {
		panic(fmt.Sprintf("aio.LimitedCounter: size should be positive number"))
	}

	lp := new(LimitedCounter)
	lp.counter = make(chan int, size)

	return lp
}

// ApplyOne applies one counter for starting new goroutine.
func (lp *LimitedCounter) ApplyOne() {
	lp.counter <- 1
	lp.wg.Add(1)
}

// ReturnOne returns one counter for ending a goroutine.
func (lp *LimitedCounter) ReturnOne() {
	<-lp.counter
	lp.wg.Done()
}

// ApplyN applys delta number counter to start new goroutines.
func (lp *LimitedCounter) ApplyN(delta int) {
	for i := 0; i < delta; i++ {
		lp.counter <- 1
	}
	for i := 0; i > delta; i-- {
		<-lp.counter
	}

	lp.wg.Add(delta)
}

// ReturnN return delta counter for ending the same amount of goroutines.
func (lp *LimitedCounter) ReturnN(delta int) {
	if delta <= 0 {
		return
	}

	for i := 0; i < delta; i++ {
		<-lp.counter
	}

	lp.wg.Done()
}

// Await blocks until the waitgroup is zero.
func (lp *LimitedCounter) Await() {
	lp.wg.Wait()
}
