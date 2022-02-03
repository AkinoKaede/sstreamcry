package net

import (
	"fmt"
)

type AddressPort struct {
	Address Address
	Port    Port
}

func (ap AddressPort) String() string {
	return fmt.Sprintf("%s:%d", ap.Address, ap.Port)
}

type Destination struct {
	Address Address
	Port    Port
	Network Network
}

func (d Destination) AddressPort() AddressPort {
	return AddressPort{
		Address: d.Address,
		Port:    d.Port,
	}
}

func TCPDestination(address Address, port Port) Destination {
	return Destination{
		Network: Network_TCP,
		Address: address,
		Port:    port,
	}
}

func UDPDestination(address Address, port Port) Destination {
	return Destination{
		Network: Network_UDP,
		Address: address,
		Port:    port,
	}
}
