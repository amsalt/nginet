package message

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/amsalt/log"
)

var (
	// ErrMessagePointerRequired represents an error require a pointer.
	ErrMessagePointerRequired = errors.New("message pointer required")
	// ErrUnnamedMessage represent an error that register message with an empty name.
	ErrUnnamedMessage = errors.New("Unnamed message")
)

// Register represents a message register.
type register struct {
	msgNameInfo map[string]Meta
	msgIDInfo   map[interface{}]Meta
}

// NewRegister creates a new register instance and return the pointer.
func NewRegister() Register {
	r := new(register)
	r.msgNameInfo = make(map[string]Meta)
	r.msgIDInfo = make(map[interface{}]Meta)

	return r
}

// GetMetaByID return the MetaData of registered protocol message by id.
func (r *register) GetMetaByID(id interface{}) Meta {
	return r.msgIDInfo[fmt.Sprintf("%v", id)]
}

// GetMetaByMsg return the MetaData of registered protocol message by message's pointer.
func (r *register) GetMetaByMsg(msg interface{}) Meta {
	mType := reflect.TypeOf(msg)
	if mType == nil || mType.Kind() != reflect.Ptr {
		panic(ErrMessagePointerRequired)
	}

	msgName := mType.Elem().Name()
	return r.msgNameInfo[msgName]
}

// RegisterMsg registers protocol messsage by message,
// the id of the message equals the name of the message.
// return MetaData info.
func (r *register) RegisterMsg(msg interface{}) (meta Meta) {
	m := r.registerMsgByName(msg)
	r.msgIDInfo[m.msgName] = m
	return m
}

// RegisterMsgByID register msg by msg ID.
func (r *register) RegisterMsgByID(assignID interface{}, msg interface{}) Meta {
	metaData := r.registerMsgByName(msg)
	metaData.msgID = assignID
	r.msgIDInfo[fmt.Sprintf("%v", assignID)] = metaData

	log.Debugf("register metaData %+v by id: %+v", metaData, metaData.msgID)
	return metaData
}

func (r *register) registerMsgByName(msg interface{}) (meta *metaData) {
	mType := reflect.TypeOf(msg)
	if mType == nil || mType.Kind() != reflect.Ptr {
		panic(ErrMessagePointerRequired)
	}

	msgName := mType.Elem().Name()
	if msgName == "" {
		panic(ErrUnnamedMessage)
	}

	if _, ok := r.msgNameInfo[msgName]; ok {
		log.Warningf("message %v is already registered", msgName)
	}

	meta = newMetaData(msgName, msgName, mType)
	r.msgNameInfo[msgName] = meta
	return
}
