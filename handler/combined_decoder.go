package handler

import (
	"encoding/binary"
	"fmt"

	"github.com/amsalt/log"
	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
)

type CombinedDecoder struct {
	*core.DefaultInboundHandler

	idParser            *IDParser
	messageDeserializer *MessageDeserializer
}

func NewCombinedDecoder(messageDeserializer *MessageDeserializer, idParser *IDParser) *CombinedDecoder {
	cd := &CombinedDecoder{DefaultInboundHandler: core.NewDefaultInboundHandler()}
	cd.idParser = idParser
	cd.messageDeserializer = messageDeserializer

	return cd
}

func (cd *CombinedDecoder) OnRead(ctx *core.ChannelContext, msg interface{}) {
	log.Debugf("CombinedDecoder.OnRead msg %T", msg)
	if data, ok := msg.(bytes.ReadOnlyBuffer); ok {
		log.Debugf("CombinedDecoder.OnRead ok")
		isCombined, err := cd.checkFlag(data, ExtraMsgFlagLength)
		if err == nil && isCombined {

			data.Discard(ExtraMsgFlagLength)
			lenBuf, err := data.Read(0, ExtraMsgLength)
			if err != nil {
				ctx.FireError(fmt.Errorf("CombinedDecoder.OnRead bad message"))
				return
			}

			msgLen := binary.BigEndian.Uint32(lenBuf)
			extraBytes, err := data.Read(0, int(msgLen))
			extraBuf := bytes.NewReadOnlyBufferWithBytes(extraBytes)
			id, extraMsg, err := cd.idParser.DecodeID(extraBuf)

			if err == nil {
				extra, err := cd.messageDeserializer.DecodePayload(id, extraMsg)
				if err == nil {
					var output []interface{}
					output = append(output, msg, extra)
					ctx.FireRead(output)
				} else {
					log.Errorf("MessageDecoder.OnRead failed: %+v", err)
				}
			} else {
				log.Errorf("CombinedDecoder.OnRead failed: %+v", err)
			}

		} else {
			ctx.FireRead(msg)
		}
	} else {
		ctx.FireError(fmt.Errorf("CombinedDecoder.OnRead msg not bytes.ReadOnlyBuffer"))
	}
}

func (cd *CombinedDecoder) checkFlag(data bytes.ReadOnlyBuffer, flagLen int) (bool, error) {
	flag, err := data.Seek(flagLen)
	if err != nil {
		return false, err
	}

	return binary.BigEndian.Uint32(flag) == ExtraMsgFlag, nil
}
