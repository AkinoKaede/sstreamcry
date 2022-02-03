package shadowsocks

import (
	"github.com/AkinoKaede/sstreamcry/common/crypto"
)

type CipherType int

const (
	CipherType_AES_128_CFB CipherType = iota
	CipherType_AES_256_CFB
)

var CipherMap = map[CipherType]Cipher{
	CipherType_AES_128_CFB: &AesCfb{KeyBytes: 16},
	CipherType_AES_256_CFB: &AesCfb{KeyBytes: 32},
}

type Cipher interface {
	KeySize() int32
	IVSize() int32
	EncodePacket(key []byte, b []byte) error
}

type AesCfb struct {
	KeyBytes int32
}

func (v *AesCfb) KeySize() int32 {
	return v.KeyBytes
}

func (v *AesCfb) IVSize() int32 {
	return 16
}

func (v *AesCfb) EncodePacket(key []byte, b []byte) error {
	iv := b[:v.IVSize()]
	stream := crypto.NewAesEncryptionStream(key, iv)
	stream.XORKeyStream(b[v.IVSize():], b[v.IVSize():])
	return nil
}
