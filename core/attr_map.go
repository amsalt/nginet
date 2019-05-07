package core

import "sync"

// AttrMap represents a map-like data structure which is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
// currently it just a wrapper of sync.Map.
type AttrMap struct {
	sync.Map
}

// NewAttrMap returns a pointer of AttrMap instance.
func NewAttrMap() *AttrMap {
	return &AttrMap{}
}

// SetValue sets a key with value.
func (attr *AttrMap) SetValue(key string, value interface{}) {
	attr.Store(key, value)
}

// SetIfAbsent Atomically sets to the given value if this value is not set and return nil.
// If it contains a value it will just return the old value and do nothing.
func (attr *AttrMap) SetIfAbsent(key string, value interface{}) interface{} {
	retValue, loaded := attr.LoadOrStore(key, value)
	if loaded {
		return retValue
	}
	return nil
}

// Value returns the value by key.
func (attr *AttrMap) Value(key string) interface{} {
	value, _ := attr.Load(key)
	return value
}

// IntValue helper method to return a int value by key.
func (attr *AttrMap) IntValue(key string) int {
	v, ok := attr.Value(key).(int)
	if ok {
		return v
	}
	return 0
}

// Int32Value helper method to return a int32 value by key.
func (attr *AttrMap) Int32Value(key string) int32 {
	v, ok := attr.Value(key).(int32)
	if ok {
		return v
	}
	return 0
}

// Uint32Value helper method to return a uint32 value by key.
func (attr *AttrMap) Uint32Value(key string) uint32 {
	v, ok := attr.Value(key).(uint32)
	if ok {
		return v
	}
	return 0
}

// Int64Value helper method to return a int64 value by key.
func (attr *AttrMap) Int64Value(key string) int64 {
	v, ok := attr.Value(key).(int64)
	if ok {
		return v
	}
	return 0
}

// Uint64Value helper method to return a uint64 value by key.
func (attr *AttrMap) Uint64Value(key string) uint64 {
	v, ok := attr.Value(key).(uint64)
	if ok {
		return v
	}
	return 0
}

// StringValue helper method to return a string value by key.
func (attr *AttrMap) StringValue(key string) string {
	v, ok := attr.Value(key).(string)
	if ok {
		return v
	}
	return ""
}
