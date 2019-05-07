package codes

import "fmt"

// Package codes defines the common error codes used by gnet.

// A Code is an unsigned 32-bit error code
type Code uint32

const (
	// OK returns when success.
	OK Code = 0

	// Failed represents a common error.
	Failed Code = 1
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case Failed:
		return "Failed"

	default:
		return fmt.Sprintf("Code (%v)", uint32(c))
	}
}
