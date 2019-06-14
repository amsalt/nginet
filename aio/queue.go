package aio

import (
	"errors"
	"sync"
)

// ErrQueueEnded represents the queue is ended, no more operation is permitted.
var ErrQueueEnded = errors.New("queue ended")

// Queue represents a thread-safe queue based on slice for extensible.
// Unlike a queue based on go channel, there is no upper limit on capacity
// of the queue if regardless of hardware limitations, so that it won't block
// when push new element into the queue.
type Queue struct {
	list  []interface{}
	mutex sync.Mutex
	cond  *sync.Cond
	wait  int
}

// NewQueue return new queue based on slice.
func NewQueue() *Queue {
	queue := &Queue{}
	queue.cond = sync.NewCond(&queue.mutex)
	return queue
}

// Push add item to queue and won't be blocked
func (queue *Queue) Push(element interface{}) {
	needSignal := false
	queue.mutex.Lock()
	queue.list = append(queue.list, element)
	if queue.wait > 0 {
		needSignal = true
	}
	queue.mutex.Unlock()

	if needSignal {
		queue.cond.Signal()
	}
}

// Poll pop the first elment
func (queue *Queue) Poll() (element interface{}) {
	if len(queue.list) > 0 {
		element = queue.list[0]
		queue.list = queue.list[1:]
	}

	return
}

// Peek same with `Poll()` return the first element but not pop.
func (queue *Queue) Peek() (element interface{}) {
	if len(queue.list) > 0 {
		element = queue.list[0]
	}

	return
}

// PollAll return a copy of all element in queue
func (queue *Queue) PollAll(copy *[]interface{}) error {
	queue.mutex.Lock()
	var err error

	for len(queue.list) == 0 {
		queue.wait++
		queue.cond.Wait()
		queue.wait--
	}

	for _, data := range queue.list {
		if data == nil {
			err = ErrQueueEnded
			break
		} else {
			*copy = append(*copy, data)
		}
	}

	queue.reset()
	queue.mutex.Unlock()

	return err
}

func (queue *Queue) reset() {
	queue.list = queue.list[:0]
}
