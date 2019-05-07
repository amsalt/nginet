package packet

import "github.com/amsalt/nginet/message"

type RawPacket struct {
	msgID   interface{}
	rawData []byte
}

func NewRawPacket(id interface{}, raw []byte) message.Packet {
	return &RawPacket{msgID: id, rawData: raw}
}

func (r *RawPacket) ID() interface{} {
	return r.msgID
}

func (r *RawPacket) Payload() interface{} {
	return r.rawData
}
