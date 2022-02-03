package internet

import (
	gonet "net"

	"github.com/AkinoKaede/sstreamcry/common/net"
)

func DialTCP(dest net.Destination) (gonet.Conn, error) {
	return gonet.Dial("tcp", dest.StringWithoutNetwork())
}
