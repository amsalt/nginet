package message

import (
	"fmt"

	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/safe"
)

// MsgID2ProcessorMap is an map which holds the mapping relation of message ID and its processor.
type MsgID2ProcessorMap map[interface{}]*Processor
type ProcessorFunc func(ctx *core.ChannelContext, msg interface{}, args ...interface{})

// Processor represents a wrapper of the msgID with processor callback.
type Processor struct {
	msgID interface{}
	cb    ProcessorFunc
}

func newProcessor(msgID interface{}, h ProcessorFunc) *Processor {
	return &Processor{msgID, h}
}

// Call calls the handler function with ctx and msg as parameters.
// Not safe
func (p *Processor) Call(ctx *core.ChannelContext, msg interface{}, args ...interface{}) {
	p.cb(ctx, msg, args...)
}

func (p *Processor) SafeCall(ctx *core.ChannelContext, msg interface{}, args ...interface{}) {
	safe.Call(func() {
		p.Call(ctx, msg, args...)
	})

}

type processorMgr struct {
	processors MsgID2ProcessorMap
	register   Register
}

// NewProcessorMgr creates a new ProcessorMgr instance and returns the pointer.
func NewProcessorMgr(register Register) ProcessorMgr {
	pm := new(processorMgr)
	pm.processors = make(MsgID2ProcessorMap)
	pm.register = register
	return pm
}

func (pm *processorMgr) RegisterProcessor(msg interface{}, h ProcessorFunc) error {
	meta := pm.register.GetMetaByMsg(msg)
	if meta == nil {
		return fmt.Errorf("register handler failed for message:%+v not registered", msg)
	}
	msgID := meta.ID()
	return pm.RegisterProcessorByID(msgID, h)
}

func (pm *processorMgr) RegisterProcessorByID(msgID interface{}, hf ProcessorFunc) error {
	if pm.processors[fmt.Sprintf("%v", msgID)] != nil {
		return fmt.Errorf("register handler failed for message id:%+v already registered", msgID)
	}

	p := newProcessor(fmt.Sprintf("%v", msgID), hf)
	pm.processors[fmt.Sprintf("%v", msgID)] = p
	return nil
}

func (pm *processorMgr) GetProcessorByID(msgID interface{}) *Processor {
	return pm.processors[fmt.Sprintf("%v", msgID)]
}
