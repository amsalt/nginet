package handler

import (
	"crypto/md5"
	"crypto/rc4"
	"fmt"

	"github.com/amsalt/nginet/bytes"
	"github.com/amsalt/nginet/core"
)

type Rc4Cipher struct {
	*core.DefaultInboundHandler
	*core.DefaultOutboundHandler

	cipher *rc4.Cipher
}

func NewRc4Cipher(key string) *Rc4Cipher {
	rc := &Rc4Cipher{
		DefaultInboundHandler:  core.NewDefaultInboundHandler(),
		DefaultOutboundHandler: core.NewDefaultOutboundHandler(),
	}
	key16Bits := md5.Sum([]byte(key))

	var err error
	rc.cipher, err = rc4.NewCipher([]byte(key16Bits[:]))
	if err != nil {
		panic(fmt.Errorf("create rc4 cipher failed for %+v", err))
	}
	return rc
}

func (rc *Rc4Cipher) OnRead(ctx *core.ChannelContext, msg interface{}) {
	if buff, ok := msg.(bytes.ReadOnlyBuffer); ok {
		rc.decrypt(buff.Bytes())
		ctx.FireRead(msg)
	} else {
		ctx.FireError(fmt.Errorf("Rc4Cipher.OnRead msg not bytes.ReadOnlyBuffer"))
	}
}

func (rc *Rc4Cipher) OnWrite(ctx *core.ChannelContext, msg interface{}) {
	if buff, ok := msg.([]byte); ok {
		rc.encrypt(buff)
	} else if buff, ok := msg.(bytes.WriteOnlyBuffer); ok {
		rc.encrypt(buff.Bytes())
	}

	ctx.FireWrite(msg)
}

func (rc *Rc4Cipher) encrypt(raw []byte) {
	rc.cipher.XORKeyStream(raw, raw)
}

func (rc *Rc4Cipher) decrypt(raw []byte) {
	rc.encrypt(raw)
}
