package net

import (
	"bytes"
	"net"
	"strings"
)

type AddressFamily byte

const (
	AddressFamilyIPv4 AddressFamily = iota
	AddressFamilyIPv6
	AddressFamilyDomain
)

type Address interface {
	IP() net.IP
	Domain() string
	Family() AddressFamily
	String() string
}

type ipv4Address [4]byte

func (a ipv4Address) IP() net.IP {
	return net.IP(a[:])
}

func (ipv4Address) Domain() string {
	panic("Calling Domain() on an IPv4Address.")
}

func (ipv4Address) Family() AddressFamily {
	return AddressFamilyIPv4
}

func (a ipv4Address) String() string {
	return a.IP().String()
}

type ipv6Address [16]byte

func (a ipv6Address) IP() net.IP {
	return net.IP(a[:])
}

func (ipv6Address) Domain() string {
	panic("Calling Domain() on an IPv6Address.")
}

func (ipv6Address) Family() AddressFamily {
	return AddressFamilyIPv6
}

func (a ipv6Address) String() string {
	return "[" + a.IP().String() + "]"
}

type domainAddress string

func (domainAddress) IP() net.IP {
	panic("Calling IP() on a DomainAddress.")
}

func (a domainAddress) Domain() string {
	return string(a)
}

func (domainAddress) Family() AddressFamily {
	return AddressFamilyDomain
}

func (a domainAddress) String() string {
	return a.Domain()
}

var bytes0 = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func IPAddress(ip []byte) Address {
	switch len(ip) {
	case net.IPv4len:
		var addr ipv4Address = [4]byte{ip[0], ip[1], ip[2], ip[3]}
		return addr
	case net.IPv6len:
		if bytes.Equal(ip[:10], bytes0) && ip[10] == 0xff && ip[11] == 0xff {
			return IPAddress(ip[12:16])
		}
		var addr ipv6Address = [16]byte{
			ip[0], ip[1], ip[2], ip[3],
			ip[4], ip[5], ip[6], ip[7],
			ip[8], ip[9], ip[10], ip[11],
			ip[12], ip[13], ip[14], ip[15],
		}
		return addr
	default:
		return nil
	}
}

func DomainAddress(domain string) Address {
	return domainAddress(domain)
}

func isAlphaNum(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func ParseAddress(addr string) Address {
	// Handle IPv6 address in form as "[2001:4860:0:2001::68]"
	lenAddr := len(addr)
	if lenAddr > 0 && addr[0] == '[' && addr[lenAddr-1] == ']' {
		addr = addr[1 : lenAddr-1]
		lenAddr -= 2
	}

	if lenAddr > 0 && (!isAlphaNum(addr[0]) || !isAlphaNum(addr[len(addr)-1])) {
		addr = strings.TrimSpace(addr)
	}

	ip := net.ParseIP(addr)
	if ip != nil {
		return IPAddress(ip)
	}
	return DomainAddress(addr)
}
