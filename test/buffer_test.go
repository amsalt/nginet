package test

import (
	"testing"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/log"
)

func TestWriteOnlyBuffer(t *testing.T) {
	wb := bytes.NewWriteOnlyBuffer(4)
	wb.WriteHeader([]byte{'3', '4'})
	log.Infof("buffer bytes: %+v", wb.Bytes())
	b, _ := wb.TakeFreeHeader(2)
	log.Infof("buffer bytes: %+v", wb.Bytes())
	b[0] = '1'
	b[1] = '2'
	log.Infof("buffer bytes: %+v", wb.Bytes())
	wb.WriteTail([]byte{'5', '6'})
	log.Infof("buffer bytes: %+v", wb.Bytes())

	s := bytes.NewWriteOnlyBufferWithBytes(4, wb.Bytes())
	s.WriteHeader([]byte("12"))
	log.Infof("create new write buf: %+v\n", s.Bytes())

	wb.Reset()
	log.Infof("reset result: %+v", wb.Bytes())

}

func TestReadOnlyBuffer(t *testing.T) {
	rb := bytes.NewReadOnlyBuffer(6)
	b := rb.FreeBytes()
	log.Info(len(b))
	b[0] = '0'
	b[1] = '1'
	b[2] = '2'
	b[3] = '3'
	log.Infof("ReadOnlyBuffer bytes: %+v, %+v", rb.Bytes(), rb.Len())
	buf, _ := rb.Seek(3)
	log.Infof("ReadOnlyBuffer seek bytes: %+v", buf)
}
