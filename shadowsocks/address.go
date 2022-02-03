package shadowsocks

import (
	"bytes"
	"encoding/binary"

	"github.com/AkinoKaede/sstreamcry/common/net"
)

func ParseDestination(dest net.Destination) []byte {
	buf := bytes.NewBuffer(nil)
	switch dest.Address.Family() {
	case net.AddressFamilyIPv4:
		buf.WriteByte(0x01)
		buf.Write([]byte(dest.Address.IP()))
	case net.AddressFamilyIPv6:
		buf.WriteByte(0x04)
		buf.Write([]byte(dest.Address.IP()))
	case net.AddressFamilyDomain:
		domain := dest.Address.Domain()
		buf.Write([]byte{0x03, byte(len(domain))})
		buf.WriteString(domain)
	}

	binary.Write(buf, binary.BigEndian, dest.Port)

	return buf.Bytes()
}
