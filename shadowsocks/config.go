package shadowsocks

type Account struct {
	Key    []byte
	Cipher Cipher
}

func CreateAccount(password, mothod string) (*Account, error) {
	cipherType, err := CipherFromString(mothod)
	if err != nil {
		return nil, err
	}

	cipher := CipherMap[cipherType]
	key := passwordToCipherKey([]byte(password), cipher.KeySize())
	return &Account{
		Key:    key,
		Cipher: cipher,
	}, nil
}
