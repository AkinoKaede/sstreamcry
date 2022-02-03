package net

import (
	"fmt"
)

type Destination struct {
	Address Address
	Port    Port
	Network Network
}

func (d Destination) StringWithoutNetwork() string {
	return fmt.Sprintf("%s:%d", d.Address, d.Port)
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
