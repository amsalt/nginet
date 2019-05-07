package idparser

import "errors"

// packae idparser implements all kinds of the packet ID parsers.
// It contains uint16 and uint32 id parser currently.

var (
	// ErrMsgIDNotCompleted represents an error that read id from bytes but it's not long enough.
	ErrMsgIDNotCompleted = errors.New("idparser: msg id not completed")

	// ErrMsgIDConvertIntFailed represents an error that failed to convert packet id to int type.
	ErrMsgIDConvertIntFailed = errors.New("idparser: msg id  convert int failed")
)

const (
	// U16IDLength the length of a byte array of type uint16
	U16IDLength = 2

	// U32IDLength the length of a byte array of type uint32
	U32IDLength = 4
)
