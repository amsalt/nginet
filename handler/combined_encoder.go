package handler

import (
	"encoding/binary"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
	"github.com/amsalt/nginet/message/packet"
)

type CombinedEncoder struct {
	*core.DefaultOutboundHandler

	idParser          *IDParser
	messageSerializer *MessageSerializer
}

func NewCombinedEncoder(messageSerializer *MessageSerializer, idParser *IDParser) *CombinedEncoder {
	ce := &CombinedEncoder{DefaultOutboundHandler: core.NewDefaultOutboundHandler()}
	ce.idParser = idParser
	ce.messageSerializer = messageSerializer

	return ce
}

func (ce *CombinedEncoder) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	log.Debugf("CombinedEncoder.OnWrite msg: %T", msg)
	if combined, ok := msg.([]interface{}); ok { // support combined message.
		log.Debugf("combined %+v ok? %+v, len: %+v", combined, ok, len(combined))
		if len(combined) >= 2 {
			log.Debugf("combined[0] %T, combined[1] %T", combined[0], combined[1])
			rawPacket, ok := combined[0].(*packet.RawPacket)
			if !ok {
				ctx.FireWrite(msg)
				return
			}

			idBuf, _, err := ce.idParser.EncodeID(rawPacket)
			log.Debugf("encode raw packet id: %+v", idBuf)
			if err != nil {
				ctx.FireError(err)
				return
			}
			payload, ok := rawPacket.Payload().([]byte)
			if !ok {
				ctx.FireWrite(msg)
				return
			}

			idBuf.WriteTail(payload)

			extra := combined[1]
			extraBuf, err := ce.encode(extra)
			log.Debugf("encode extraBuf: %+v, err: %+v", extraBuf, err)
			if err != nil {
				ctx.FireError(err)
			} else {
				// use uint16 as length field.
				length := make([]byte, ExtraMsgLength)
				binary.BigEndian.PutUint32(length, uint32(extraBuf.Len()))
				extraBuf.WriteHeader(length)

				// use 4 bytes as the flag of extra fields
				// if not the flag, represents no extra fields.
				extraFlag := make([]byte, ExtraMsgFlagLength)
				binary.BigEndian.PutUint32(extraFlag, ExtraMsgFlag)
				extraBuf.WriteHeader(extraFlag)

				extraBuf.WriteTail(idBuf.Bytes())
				log.Errorf("CombinedEncoder onwrite: %+v", extraBuf.Bytes())
				ctx.FireWrite(extraBuf)
			}
		} else {
			ctx.FireWrite(msg)
		}

	} else {
		ctx.FireWrite(msg)
	}
}

func (ce *CombinedEncoder) encode(msg interface{}) (bytes.WriteOnlyBuffer, error) {
	buf, id, err := ce.idParser.EncodeID(msg)
	if err != nil {
		return nil, err
	}

	err = ce.messageSerializer.EncodePayload(buf, msg, id)
	return buf, err
}
