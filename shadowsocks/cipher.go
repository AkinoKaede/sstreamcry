package shadowsocks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rc4"
	"strings"

	"github.com/AkinoKaede/sstreamcry/common"
	"github.com/aead/chacha20"
	"github.com/aead/chacha20/chacha"
	"github.com/kierdavis/cfb8"
)

type CipherType int

const (
	CipherType_UNKNOWN CipherType = iota
	CipherType_AES_128_CTR
	CipherType_AES_192_CTR
	CipherType_AES_256_CTR
	CipherType_AES_128_CFB
	CipherType_AES_192_CFB
	CipherType_AES_256_CFB
	CipherType_AES_128_CFB8
	CipherType_AES_192_CFB8
	CipherType_AES_256_CFB8
	CipherType_CHACHA20
	CipherType_CHACHA20_IETF
	CipherType_XCHACHA20
	CipherType_RC4
	CipherType_RC4_MD5
)

var CipherMap = map[CipherType]Cipher{
	CipherType_AES_128_CTR: &StreamCipher{
		KeyBytes:       16,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
	},
	CipherType_AES_192_CTR: &StreamCipher{
		KeyBytes:       24,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
	},
	CipherType_AES_256_CTR: &StreamCipher{
		KeyBytes:       32,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
	},
	CipherType_AES_128_CFB: &StreamCipher{
		KeyBytes:       16,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewCFBEncrypter),
	},
	CipherType_AES_192_CFB: &StreamCipher{
		KeyBytes:       24,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewCFBEncrypter),
	},
	CipherType_AES_256_CFB: &StreamCipher{
		KeyBytes:       32,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewCFBEncrypter),
	},
	CipherType_AES_128_CFB8: &StreamCipher{
		KeyBytes:       16,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
	},
	CipherType_AES_192_CFB8: &StreamCipher{
		KeyBytes:       24,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
	},
	CipherType_AES_256_CFB8: &StreamCipher{
		KeyBytes:       32,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
	},
	CipherType_CHACHA20: &StreamCipher{
		KeyBytes: chacha.KeySize,
		IVBytes:  chacha.NonceSize,
		EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
			return chacha20.NewCipher(iv, key)
		},
	},
	CipherType_CHACHA20_IETF: &StreamCipher{
		KeyBytes: chacha.KeySize,
		IVBytes:  chacha.INonceSize,
		EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
			return chacha20.NewCipher(iv, key)
		},
	},
	CipherType_XCHACHA20: &StreamCipher{
		KeyBytes: chacha.KeySize,
		IVBytes:  chacha.INonceSize,
		EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
			return chacha20.NewCipher(iv, key)
		},
	},
	CipherType_RC4: &StreamCipher{
		KeyBytes: 16,
		IVBytes:  16,
		EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
			return rc4.NewCipher(key)
		},
	},
	CipherType_RC4_MD5: &StreamCipher{
		KeyBytes: 16,
		IVBytes:  16,
		EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
			h := md5.New()
			h.Write(key)
			h.Write(iv)
			return rc4.NewCipher(h.Sum(nil))
		},
	},
}

type Cipher interface {
	KeySize() int32
	IVSize() int32
	EncodePacket(key []byte, b []byte) error
}

func blockStream(blockCreator func(key []byte) (cipher.Block, error), streamCreator func(block cipher.Block, iv []byte) cipher.Stream) func([]byte, []byte) (cipher.Stream, error) {
	return func(key []byte, iv []byte) (cipher.Stream, error) {
		block, err := blockCreator(key)
		if err != nil {
			return nil, err
		}
		return streamCreator(block, iv), err
	}
}

type StreamCipher struct {
	KeyBytes       int32
	IVBytes        int32
	EncryptCreator func(key []byte, iv []byte) (cipher.Stream, error)
}

func (v *StreamCipher) KeySize() int32 {
	return v.KeyBytes
}

func (v *StreamCipher) IVSize() int32 {
	return v.IVBytes
}

func (v *StreamCipher) EncodePacket(key []byte, b []byte) error {
	iv := b[:v.IVSize()]
	stream, err := v.EncryptCreator(key, iv)
	if err != nil {
		return err
	}
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
	case "aes-128-ctr":
		return CipherType_AES_128_CTR
	case "aes-192-ctr":
		return CipherType_AES_192_CTR
	case "aes-256-ctr":
		return CipherType_AES_256_CTR
	case "aes-128-cfb":
		return CipherType_AES_128_CFB
	case "aes-192-cfb":
		return CipherType_AES_192_CFB
	case "aes-256-cfb":
		return CipherType_AES_256_CFB
	case "aes-128-cfb8":
		return CipherType_AES_128_CFB8
	case "aes-192-cfb8":
		return CipherType_AES_192_CFB8
	case "aes-256-cfb8":
		return CipherType_AES_256_CFB8
	case "chacha20":
		return CipherType_CHACHA20
	case "chacha20-ietf":
		return CipherType_CHACHA20_IETF
	case "xchacha20":
		return CipherType_XCHACHA20
	case "rc4":
		return CipherType_RC4
	case "rc4-md5":
		return CipherType_RC4_MD5
	default:
		return CipherType_UNKNOWN
	}
}
