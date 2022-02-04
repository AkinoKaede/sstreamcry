package shadowsocks

import "errors"

type Account struct {
	Key    []byte
	Cipher Cipher
}

func CreateAccount(password, mothod string) (*Account, error) {
	cipherType := CipherFromString(mothod)
	if cipherType == CipherType_UNKNOWN {
		return nil, errors.New("unknown cipher method: " + mothod)
	}

	cipher := CipherMap[cipherType]
	key := passwordToCipherKey([]byte(password), cipher.KeySize())
	return &Account{
		Key:    key,
		Cipher: cipher,
	}, nil
}
