package idparser

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/message"
)

type uint32ID struct {
	byteorder binary.ByteOrder
}

func NewUint32ID() message.PacketIDParser {
	return &uint32ID{byteorder: binary.BigEndian}
}

func (i *uint32ID) IDLen() int {
	return U32IDLength
}

func (i *uint32ID) SetByteOrder(order binary.ByteOrder) message.PacketIDParser {
	i.byteorder = order
	return i
}

func (i *uint32ID) Encode(packetID interface{}, wob bytes.WriteOnlyBuffer) error {
	sid := fmt.Sprintf("%v", packetID)
	nid, err := strconv.Atoi(sid)
	if err != nil {
		err = ErrMsgIDConvertIntFailed
		return err
	}

	b, err := wob.TakeFreeHeader(U32IDLength)
	if err != nil {
		return err
	}
	i.byteorder.PutUint32(b, uint32(nid))
	return err
}

func (i *uint32ID) Decode(rob bytes.ReadOnlyBuffer) (id interface{}, err error) {
	if rob.Len() < U32IDLength {
		return 0, ErrMsgIDNotCompleted // data not complete, ignore data.
	}

	idBytes, err := rob.Read(0, U32IDLength)
	if err != nil {
		return nil, err
	}
	id = i.byteorder.Uint32(idBytes)
	return
}
