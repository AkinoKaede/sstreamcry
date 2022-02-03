package internet

import (
	gonet "net"

	"github.com/AkinoKaede/sstreamcry/common/net"
)

func Dial(dest net.Destination) (gonet.Conn, error) {
	var netStr string
	switch dest.Network {
	case net.Network_TCP:
		netStr = "tcp"
	case net.Network_UDP:
		netStr = "udp"
	}

	return gonet.Dial(netStr, dest.AddressPort().String())
}
