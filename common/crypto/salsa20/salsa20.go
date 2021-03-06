package salsa20

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"

	"github.com/AkinoKaede/sstreamcry/common"
	"golang.org/x/crypto/salsa20/salsa"
)

type Cipher struct {
	nonce   []byte
	key     [32]byte
	counter uint64
}

func (c *Cipher) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		common.Must(errors.New("dst buffer is to small"))
	}
	paddingLength := int(c.counter % 64)
	buf := make([]byte, len(src)+paddingLength)

	var subNonce [16]byte
	copy(subNonce[:], c.nonce)
	binary.LittleEndian.PutUint64(subNonce[8:], c.counter/64)

	copy(buf[paddingLength:], src)
	salsa.XORKeyStream(buf, buf, &subNonce, &c.key)
	copy(dst, buf[paddingLength:])

	c.counter += uint64(len(src))
}

func New(key, nonce []byte) (cipher.Stream, error) {
	var fixedSizedKey [32]byte
	if len(key) != 32 {
		return nil, errors.New("key size must be 32")
	}
	copy(fixedSizedKey[:], key)
	return &Cipher{
		key:   fixedSizedKey,
		nonce: nonce,
	}, nil
}
