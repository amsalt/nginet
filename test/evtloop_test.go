package test

import (
	"testing"
	"time"

	"github.com/amsalt/nginet/aio"
	"github.com/amsalt/log"
)

func TestEventLoop(t *testing.T) {
	evtloop := aio.NewEventLoop()
	evtloop.Start()

	evtloop.ScheduleAtFixedRate(func() {
		log.Info(time.Now().Unix())
	}, time.Second)

	evtloop.ScheduleAtFixedDelay(func() {
		log.Infof("delay 3s test")
	}, time.Second*3)

	time.Sleep(time.Second * 5)
}
