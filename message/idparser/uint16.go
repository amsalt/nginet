package idparser

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/message"
	"github.com/amsalt/log"
)

type uint16ID struct {
	byteorder binary.ByteOrder
}

func NewUint16ID() message.PacketIDParser {
	return &uint16ID{byteorder: binary.BigEndian}
}

func (i *uint16ID) IDLen() int {
	return U16IDLength
}

func (i *uint16ID) SetByteOrder(order binary.ByteOrder) message.PacketIDParser {
	i.byteorder = order
	return i
}

func (i *uint16ID) Encode(packetID interface{}, wob bytes.WriteOnlyBuffer) error {
	sid := fmt.Sprintf("%v", packetID)
	nid, err := strconv.Atoi(sid)
	if err != nil {
		err = ErrMsgIDConvertIntFailed
		return err
	}

	b, err := wob.TakeFreeHeader(U16IDLength)
	if err != nil {
		return err
	}
	i.byteorder.PutUint16(b, uint16(nid))
	return err
}

func (i *uint16ID) Decode(rob bytes.ReadOnlyBuffer) (id interface{}, err error) {
	log.Debugf("uint16ID.Decode rob.Len: %+v, U16IDLength: %+v", rob.Len(), U16IDLength)
	if rob.Len() < U16IDLength {
		return 0, ErrMsgIDNotCompleted // data not complete, ignore data.
	}

	idBytes, err := rob.Read(0, U16IDLength)
	if err != nil {
		return nil, err
	}

	id = i.byteorder.Uint16(idBytes)
	return
}
