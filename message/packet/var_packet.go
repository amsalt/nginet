package packet

import "github.com/amsalt/nginet/message"

type VarPacket struct {
	id      interface{}
	payload interface{}
}

func NewVarPacket(packetID, payload interface{}) message.Packet {
	return &VarPacket{id: packetID, payload: payload}
}

func (v *VarPacket) ID() interface{} {
	return v.id
}

func (v *VarPacket) Payload() interface{} {
	return v.payload
}
