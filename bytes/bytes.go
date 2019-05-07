package bytes

import "errors"

// package bytes implements a new buffer for efficient bytes operation.
//
//  Firstly, it's a buffer of bytes that avoids memory allocation frequently - size customizable pre-allocated memory
// 	Secondly, avoids unnecesssary copy for read&write data.
//  Thirdly, it supports Header & Tail operation - when write data to connection, write to header will reduce copy.

var (
	// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
	ErrTooLarge = errors.New("bytes.WriteOnlyBuffer: too large")
)

func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}
