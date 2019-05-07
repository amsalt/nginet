package aio

import (
	"time"

	"github.com/amsalt/nginet/exception"
	"github.com/amsalt/log"
)

// EventLoop schedules event in one goroutine.
// Supports schedule at once or schedule at fixed rate.
type EventLoop struct {
	evtQueue *Queue // stores all event to process.
}

// NewEventLoop creates a new eventloop instance.
func NewEventLoop() *EventLoop {
	evtloop := new(EventLoop)
	evtloop.evtQueue = NewQueue()

	return evtloop
}

// Start starts event loop in a new goroutine.
func (el *EventLoop) Start() {
	go el.loop()
}

// Stop stops the event loop
func (el *EventLoop) Stop() {
	el.evtQueue.Push(nil)
}

// Execute put event to queue and the event will be processed in event loop.
func (el *EventLoop) Execute(task func()) {
	log.Debugf("EventLoop Schedule task: %+v", task)
	if task != nil { // attention
		el.evtQueue.Push(task)
	}
}

// ScheduleAtFixedRate schedules event at fixed rate
// for example:
// 	ScheduleAtFixedRate(func(){print("hello")}, time.Milliseconds*5) means print "hello" at rate of 5ms
func (el *EventLoop) ScheduleAtFixedRate(task func(), period time.Duration) {
	te := &tickevent{period: period, t: task, el: el}
	te.tick()
}

// ScheduleAtFixedDelay schedule event after delay period
// for example:
// 		ScheduleAtFixedDelay(func(){print("hello")}, time.Milliseconds*5) means print "hello" once after 5ms
func (el *EventLoop) ScheduleAtFixedDelay(task func(), period time.Duration) {
	time.AfterFunc(period, func() {
		el.Execute(task)
	})
}

func (el *EventLoop) loop() {
	var copyList []interface{}

	for {
		copyList = copyList[:0]
		err := el.evtQueue.PollAll(&copyList)

		for _, task := range copyList {
			switch t := task.(type) {
			case func():
				exception.Safecall(t)
			case nil:
				// no more event
				break
			default:
				log.Debugf("unrecognized type(%T) for Scheduler ", t)
			}
		}

		if err != nil { // stoped.
			goto Stop
		}
	}
Stop:
}

type tickevent struct {
	period time.Duration
	t      func()
	el     *EventLoop
}

func (te *tickevent) tick() {
	time.AfterFunc(te.period, func() {
		te.el.Execute(te.t) // exec event
		te.tick()           // recursive schedule new loop event
	})
}
