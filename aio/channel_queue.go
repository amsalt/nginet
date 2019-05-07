package aio

import (
	"errors"
	"time"
)

// ChanQueue represents a queue based on go chan which limits the capacity.
// Just used for benchmark.
type ChanQueue struct {
	items   chan interface{}
	timeout int
}

// NewChanQueue creates Queue based channel.
func NewChanQueue(size int, timeout int) *ChanQueue {
	queue := new(ChanQueue)
	queue.items = make(chan interface{}, size)
	queue.timeout = timeout
	return queue
}

// RawItems return raw chan
func (q *ChanQueue) RawItems() chan interface{} {
	return q.items
}

// Push add new item to chan, if set timeout, will check timeout,otherwise will block
func (q *ChanQueue) Push(item interface{}) error {
	if q.timeout <= 0 {
		return q.push(item)
	}
	return q.pushWithTimeout(item, q.timeout)
}

// Get return one item in chan
func (q *ChanQueue) Get() (interface{}, error) {
	if q.timeout <= 0 {
		return q.get()
	}
	return q.getWithTimeout(q.timeout)
}

func (q *ChanQueue) push(item interface{}) error {
	select {
	case q.items <- item:
		return nil
		// default:
		// 	return errors.New("Queue full")
	}
}

func (q *ChanQueue) pushWithTimeout(item interface{}, timeoutSecs int) error {
	select {
	case q.items <- item:
		return nil
	case <-time.After(time.Duration(timeoutSecs) * time.Second):
		return errors.New("Queue full, wait timeout")
	}
}

func (q *ChanQueue) get() (interface{}, error) {
	var item interface{}
	select {
	case item = <-q.items:
		return item, nil
	default:
		return nil, errors.New("Queue empty")
	}
}

func (q *ChanQueue) getWithTimeout(timeoutMilliSec int) (interface{}, error) {
	var item interface{}
	select {
	case item = <-q.items:
		return item, nil
	case <-time.After(time.Duration(timeoutMilliSec) * time.Millisecond):
		return nil, errors.New("Queue empty, wait timeout")
	}
}
