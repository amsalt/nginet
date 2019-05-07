package bytes

import (
	"errors"
	"io"
)

var (
	// ErrNoEnoughData represents no enough data to return in a buffer.
	ErrNoEnoughData = errors.New("bytes.ReadOnlyBuffer: no engough data")
)

// ReadOnlyBuffer defines a buffer only used for reading data.
type ReadOnlyBuffer interface {
	// Read reads n bytes start with offset from buffer.
	Read(offset, n int) ([]byte, error)

	// ReadFrom reads reader's data to buffer
	ReadFrom(reader io.Reader) (int, error)

	// Seek detects bytes of n from buffer.
	Seek(n int) ([]byte, error)

	// Len returns the number of bytes of the unread portion of the buffer;
	// b.Len() == len(b.Bytes()).
	Len() int

	// FreeBytes returns all the free space in the buffer and won't mark the slice as read.
	FreeBytes() []byte

	// Discard marks the data of n from off as readed.
	Discard(n int)

	// Bytes returns a slice of length len(b.buf) holding the unread portion of the buffer.
	// The slice is valid for use only until the next buffer modification.
	// The slice aliases the buffer content at least until the next buffer modification,
	// so immediate changes to the slice will affect the result of future reads.
	Bytes() []byte

	// Reset resets the buffer to be empty,
	// but it retains the underlying storage for use by future writes.
	Reset()
}

type readOnlyBuffer struct {
	buf []byte
	off int
	end int
}

// NewReadOnlyBuffer creates new ReadOnlyBuffer instance with buffer of n.
func NewReadOnlyBuffer(n int) ReadOnlyBuffer {
	b := &readOnlyBuffer{}
	if n <= 0 {
		n = 1024 // default 1k
	}

	b.buf = make([]byte, 0, n)
	b.Reset()
	return b
}

// NewReadOnlyBufferWithBytes creates new instance of ReadOnlyBuffer with buf.
func NewReadOnlyBufferWithBytes(buf []byte) ReadOnlyBuffer {
	rob := &readOnlyBuffer{buf: buf}
	rob.off = 0
	rob.end = len(buf)
	return rob
}

func (b *readOnlyBuffer) Read(offset, n int) (buf []byte, err error) {
	if offset < 0 || offset+n > b.Len() {
		return nil, ErrNoEnoughData
	}

	buf = b.buf[b.off+offset : b.off+offset+n]
	b.off += offset + n
	return
}

func (b *readOnlyBuffer) Seek(n int) ([]byte, error) {
	if n < 0 || n > b.Len() {
		return nil, ErrNoEnoughData
	}

	return b.buf[b.off : b.off+n], nil
}

func (b *readOnlyBuffer) ensureSpace() int {
	m, ok := b.tryGrowByReslice(b.Len())
	if !ok {
		m = b.grow()
	}
	return m
}

func (b *readOnlyBuffer) grow() int {
	m := b.Len()
	c := cap(b.buf)
	buf := makeSlice(2 * c)
	copy(buf, b.buf[b.off:b.end])
	b.buf = buf
	b.end -= b.off
	b.off = 0
	return m
}

func (b *readOnlyBuffer) tryGrowByReslice(n int) (int, bool) {
	if n == 0 {
		b.buf = b.buf[:cap(b.buf)]
		return 0, true
	}

	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

func (b *readOnlyBuffer) FreeBytes() []byte {
	b.ensureSpace()
	return b.buf[b.end:]
}

func (b *readOnlyBuffer) Len() int {
	return b.end - b.off
}

func (b *readOnlyBuffer) Discard(n int) {
	if n < 0 || n > b.Len() {
		panic("bytes.ReadOnlyBuffer: Discard out of range")
	}

	b.off += n
}

func (b *readOnlyBuffer) Bytes() []byte {
	return b.buf[b.off:b.end]
}

func (b *readOnlyBuffer) ReadFrom(reader io.Reader) (int, error) {
	n, err := reader.Read(b.FreeBytes())
	if err != nil {
		return n, err
	}
	b.end += n
	return n, nil
}

func (b *readOnlyBuffer) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
	b.end = 0
}
