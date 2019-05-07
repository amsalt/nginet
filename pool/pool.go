package pool

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/amsalt/nginet/exception"
)

const (
	DefaultPoolCap                  = 1024 * 1000
	DefaultClearWorkerInterval      = 60
	DefaultCheckClearWorkerInterval = 10
)

type Pool struct {
	sync.Mutex
	sync.Cond
	sync.Once

	workers []*worker

	clearWorkerInterval int
	clearWorkerTimer    *time.Timer

	capacity int32
	running  int32
	closed   int32
}

func NewPool() *Pool {
	p := &Pool{}
	p.SetSize(DefaultPoolCap)
	p.SetClearWorkerInterval(DefaultClearWorkerInterval)
	p.startClearW()
	return p
}

func (p *Pool) SetSize(cap int32) {
	if cap == p.capacity {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(cap))

	diff := atomic.LoadInt32(&p.running) - cap
	for i := 0; i < int(diff); i++ {
		p.applyW().stop()
	}

	// otherwise: do nothing.
}

func (p *Pool) SetClearWorkerInterval(interval int) {
	p.clearWorkerInterval = interval
}

func (p *Pool) Execute(task func()) {
	w := p.applyW()
	w.runTask(task)
}

func (p *Pool) Stop() {
	p.Do(func() {
		atomic.StoreInt32(&p.closed, 1)
		p.Lock()
		p.clearWorkerTimer.Stop()
		for i, w := range p.workers {
			w.stop()
			p.workers[i] = nil
		}
		p.workers = nil
		p.Unlock()
	})
}

func (p *Pool) startClearW() {
	go p.clearW()
}

func (p *Pool) clearW() {
	var clearTask = func() {
		p.Lock()
		p.Unlock()
		index := -1
		for i, w := range p.workers {
			if p.clearWorkerInterval > 0 && time.Now().Sub(w.returnTime) >= time.Second*time.Duration(p.clearWorkerInterval) {
				p.releaseW(w)
				index = i
				p.workers[i] = nil
			} else {
				// FIFO
				break
			}
		}

		if index > -1 {
			if index >= len(p.workers)-1 {
				p.workers = p.workers[:0]
			} else {
				p.workers = p.workers[index+1:]
			}
		}
	}

	p.clearWorkerTimer = time.AfterFunc(time.Second*DefaultCheckClearWorkerInterval, func() {
		clearTask()
		p.clearW()
	})

}

func (p *Pool) releaseW(w *worker) bool {
	atomic.AddInt32(&p.running, -1)
	w.stop()
	return true
}

func (p *Pool) startNewW() *worker {
	w := &worker{pool: p}
	w.task = make(chan func())

	atomic.AddInt32(&p.running, 1)
	w.start()
	return w
}

func (p *Pool) waitFreeW() *worker {
	var w *worker
	for {
		// wait notify
		p.Wait()
		l := len(p.workers) - 1
		if l < 0 {
			continue
		}
		w = p.workers[l]
		p.workers[l] = nil
		p.workers = p.workers[:l]
		break
	}
	return w
}

func (p *Pool) applyW() *worker {
	var w *worker
	p.Lock()
	defer p.Unlock()

	freeWNum := len(p.workers)
	// free worker, no task
	if freeWNum > 0 {
		w = p.workers[freeWNum-1]
		p.workers[freeWNum-1] = nil
		p.workers = p.workers[:freeWNum-1]
	} else if p.capacity > p.running { // pool not full.
		w = p.startNewW()
	} else { // waitting for free worker.
		w = p.waitFreeW()
	}

	return w
}

func (p *Pool) returnW(w *worker) bool {
	atomic.AddInt32(&p.running, -1)
	w.returnTime = time.Now()
	p.Lock()
	p.workers = append(p.workers, w)
	p.Unlock()

	// notify free worker.
	p.Signal()
	return true
}

type worker struct {
	pool *Pool
	task chan func()

	returnTime time.Time
}

func (w *worker) start() {
	exception.TryCatch(func() {
		go w.run()
	}, func() {
		w.pool.releaseW(w)
	})
}

func (w *worker) runTask(t func()) {
	w.task <- t
}

func (w *worker) stop() {
	w.task <- nil
}

func (w *worker) run() {
	for t := range w.task {
		// end flag
		if t == nil {
			w.pool.releaseW(w)
			break
		}
		// call
		t()

		// after call, return the worker.
		if !w.pool.returnW(w) {
			break
		}
	}
}

var defaultPool *Pool

func init() {
	defaultPool = NewPool()
}

func SetSize(cap int32) {
	defaultPool.SetSize(cap)
}

func Execute(task func()) {
	defaultPool.Execute(task)
}

func Stop() {
	defaultPool.Stop()
}
