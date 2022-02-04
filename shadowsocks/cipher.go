package shadowsocks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rc4"
	"errors"
	"strings"

	"github.com/AkinoKaede/sstreamcry/common"
	"github.com/aead/chacha20"
	"github.com/aead/chacha20/chacha"
	"github.com/kierdavis/cfb8"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/cast5"
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
	CipherType_AES_128_OFB
	CipherType_AES_192_OFB
	CipherType_AES_256_OFB
	CipherType_CHACHA20
	CipherType_CHACHA20_IETF
	CipherType_XCHACHA20
	CipherType_RC4
	CipherType_RC4_MD5
	CipherType_BF_CFB
	CipherType_CAST5_CFB
	CipherType_DES_CFB
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
	CipherType_AES_128_OFB: &StreamCipher{
		KeyBytes:       16,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewOFB),
	},
	CipherType_AES_192_OFB: &StreamCipher{
		KeyBytes:       24,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewOFB),
	},
	CipherType_AES_256_OFB: &StreamCipher{
		KeyBytes:       32,
		IVBytes:        aes.BlockSize,
		EncryptCreator: blockStream(aes.NewCipher, cipher.NewOFB),
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
	CipherType_BF_CFB: &StreamCipher{
		KeyBytes:       16,
		IVBytes:        blowfish.BlockSize,
		EncryptCreator: blockStream(func(key []byte) (cipher.Block, error) { return blowfish.NewCipher(key) }, cipher.NewCFBEncrypter),
	},
	CipherType_CAST5_CFB: &StreamCipher{
		KeyBytes:       16,
		IVBytes:        cast5.BlockSize,
		EncryptCreator: blockStream(func(key []byte) (cipher.Block, error) { return cast5.NewCipher(key) }, cipher.NewCFBEncrypter),
	},
	CipherType_DES_CFB: &StreamCipher{
		KeyBytes:       8,
		IVBytes:        des.BlockSize,
		EncryptCreator: blockStream(des.NewCipher, cipher.NewCFBEncrypter),
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

func CipherFromString(c string) (CipherType, error) {
	switch strings.ToLower(c) {
	case "aes-128-ctr":
		return CipherType_AES_128_CTR, nil
	case "aes-192-ctr":
		return CipherType_AES_192_CTR, nil
	case "aes-256-ctr":
		return CipherType_AES_256_CTR, nil
	case "aes-128-cfb":
		return CipherType_AES_128_CFB, nil
	case "aes-192-cfb":
		return CipherType_AES_192_CFB, nil
	case "aes-256-cfb":
		return CipherType_AES_256_CFB, nil
	case "aes-128-cfb8":
		return CipherType_AES_128_CFB8, nil
	case "aes-192-cfb8":
		return CipherType_AES_192_CFB8, nil
	case "aes-256-cfb8":
		return CipherType_AES_256_CFB8, nil
	case "aes-128-ofb":
		return CipherType_AES_128_OFB, nil
	case "aes-192-ofb":
		return CipherType_AES_192_OFB, nil
	case "aes-256-ofb":
		return CipherType_AES_256_OFB, nil
	case "chacha20":
		return CipherType_CHACHA20, nil
	case "chacha20-ietf":
		return CipherType_CHACHA20_IETF, nil
	case "xchacha20":
		return CipherType_XCHACHA20, nil
	case "rc4":
		return CipherType_RC4, nil
	case "rc4-md5":
		return CipherType_RC4_MD5, nil
	case "bf-cfb", "blowfish-cfb":
		return CipherType_BF_CFB, nil
	case "cast5-cfb":
		return CipherType_CAST5_CFB, nil
	case "des-cfb":
		return CipherType_DES_CFB, nil
	default:
		return CipherType_UNKNOWN, errors.New("unknown cipher method: " + c)
	}
}
