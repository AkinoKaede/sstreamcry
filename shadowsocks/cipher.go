package shadowsocks

import (
	"crypto/md5"
	"log"
	"strings"

	"github.com/AkinoKaede/sstreamcry/common"
	"github.com/AkinoKaede/sstreamcry/common/crypto"
)

type Account struct {
	Key    []byte
	Cipher Cipher
}

func CreateAccount(password, mothod string) Account {
	cipherType := CipherFromString(mothod)
	cipher := CipherMap[cipherType]
	key := passwordToCipherKey([]byte(password), cipher.KeySize())
	return Account{
		Key:    key,
		Cipher: cipher,
	}
}

type CipherType int

const (
	CipherType_UNKNOWN     CipherType = iota
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

func passwordToCipherKey(password []byte, keySize int32) []byte {
	key := make([]byte, 0, keySize)

	md5Sum := md5.Sum(password)
	key = append(key, md5Sum[:]...)

	for int32(len(key)) < keySize {
		md5Hash := md5.New()
		common.Must2(md5Hash.Write(md5Sum[:]))
		common.Must2(md5Hash.Write(password))
		md5Hash.Sum(md5Sum[:0])

		key = append(key, md5Sum[:]...)
	}
	return key
}

func CipherFromString(c string) CipherType {
	switch strings.ToLower(c) {
	case "aes-128-cfb":
		return CipherType_AES_128_CFB
	case "aes-256-cfb":
		return CipherType_AES_256_CFB
	default:
		log.Fatalln("unknown cipher method:", c)
		return CipherType_UNKNOWN
	}
}
