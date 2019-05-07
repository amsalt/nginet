package bytes

import (
	"errors"
)

var (
	// ErrNoEnoughHeader represents no enough space for header to store data in a buffer.
	ErrNoEnoughHeader = errors.New("bytes.WriteOnlyBuffer: no enough header space to write")
)

// WriteOnlyBuffer defines a buffer only used for easy-write and full-read.
// For write header, when create a new WriteOnlyBuffer, a header size need to be assigned.
type WriteOnlyBuffer interface {
	// WriteHeader writes slice p to fixed length header. The return value n is the length of
	// writting in to the header. If the length of p > the free space of header,
	// WriteHeader will return with ErrNoEnoughHeader
	WriteHeader(p []byte) (n int, err error)

	// TakeFreeHeader returns the free slice of n in header and mark as used.
	// For high performance, thinking of using TakeFreeHeader instead of WriteHeader method as much as possible.
	TakeFreeHeader(n int) ([]byte, error)

	// WriteTail appends the contents of p to the tail of the buffer, growing the buffer as
	// needed. The return value n is the length of p; err is always nil. If the
	// buffer becomes too large, WriteTail will panic with ErrTooLarge.
	WriteTail(p []byte) (n int, err error)

	// FreeTail returns the free memory of tail.
	FreeTail() ([]byte, int)

	// Bytes returns a slice of length len(b.buf) holding the unread portion of the buffer.
	// The slice is valid for use only until the next buffer modification (that is,
	// only until the next call to a method like WriteTail).
	// The slice aliases the buffer content at least until the next buffer modification,
	// so immediate changes to the slice will affect the result of future reads.
	Bytes() []byte

	Len() int

	// Reset resets the buffer to be empty,
	// but it retains the underlying storage for use by future writes.
	Reset()
}

// writeOnlyBuffer is a implementation of WriteOnlyBuffer.
type writeOnlyBuffer struct {
	hlen int // the initial size of header.

	start int    // the start index of buf.
	end   int    // the end index of buf data.
	buf   []byte // contents are the bytes buf[start : len(buf)]
}

// NewWriteOnlyBuffer creates a new WriteOnlyBuffer with header length with hl, and initial tail length with tl.
func NewWriteOnlyBuffer(hl uint) WriteOnlyBuffer {
	wb := &writeOnlyBuffer{hlen: int(hl)}
	buf := make([]byte, hl)
	wb.start = int(hl)
	wb.end = wb.start
	wb.buf = buf
	return wb
}

func NewWriteOnlyBufferWithBytes(head uint, buf []byte) WriteOnlyBuffer {
	wb := &writeOnlyBuffer{}
	if head > 0 {
		newbuf := make([]byte, int(head)+len(buf), int(head)+len(buf))
		copy(newbuf[head:], buf[0:])
		wb.start = int(head)
		wb.end = wb.start + len(buf)
		wb.buf = newbuf
	} else {
		wb.start = 0
		wb.end = len(buf)
		wb.buf = buf
	}

	return wb
}

func (wb *writeOnlyBuffer) WriteHeader(p []byte) (n int, err error) {
	if len(p) > wb.start {
		return 0, ErrNoEnoughHeader
	}

	wb.start = wb.start - len(p)
	copy(wb.buf[wb.start:], p[0:])
	return len(p), nil
}

func (wb *writeOnlyBuffer) TakeFreeHeader(n int) ([]byte, error) {
	if n > wb.start {
		return nil, ErrNoEnoughHeader
	}
	wb.start -= n
	return wb.buf[wb.start : wb.start+n], nil
}

func (wb *writeOnlyBuffer) grow(n int) {
	c := cap(wb.buf)
	buf := makeSlice(2*c + n)
	copy(buf[wb.start:], wb.buf[wb.start:])
	wb.buf = buf
}

func (wb *writeOnlyBuffer) tryGrowByReslice(n int) bool {
	if n <= cap(wb.buf)-wb.end {
		wb.buf = wb.buf[:wb.end+n]
		return true
	}
	return false
}

func (wb *writeOnlyBuffer) WriteTail(p []byte) (n int, err error) {
	ok := wb.tryGrowByReslice(len(p))
	if !ok {
		wb.grow(len(p))
	}
	n = copy(wb.buf[wb.end:], p)
	wb.end += n
	return

}

func (wb *writeOnlyBuffer) FreeTail() ([]byte, int) {
	return wb.buf[wb.end:], len(wb.buf) - wb.end
}

func (wb *writeOnlyBuffer) Bytes() []byte {
	return wb.buf[wb.start:wb.end]
}

func (wb *writeOnlyBuffer) Len() int {
	return wb.end - wb.start
}

func (wb *writeOnlyBuffer) Reset() {
	wb.start = wb.hlen
	wb.end = wb.start
	wb.buf = wb.buf[:wb.start]
}
