package shadowsocks

import (
	"bytes"

	"github.com/AkinoKaede/sstreamcry/common/net"
)

func EncodeHeaderChain(destination net.Destination, account Account, times int) []byte {
	var data []byte

	for i := 0; i < times; i++ {
		data = encode(destination, data, account)
	}

	return data
}

func encode(destination net.Destination, payload []byte, account Account) []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(ParseDestination(destination))
	buf.Write(payload)

	b := buf.Bytes()
	cipher := account.Cipher
	cipher.EncodePacket(account.Key, b)

	return b
}
