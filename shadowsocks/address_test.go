package shadowsocks

import (
	"testing"

	"github.com/AkinoKaede/sstreamcry/common/net"
	"github.com/google/go-cmp/cmp"
)

func TestParseDestination(t *testing.T) {
	testCases := []struct {
		Input  net.Destination
		Output []byte
	}{
		{
			Input: net.Destination{
				Address: net.ParseAddress("1.1.1.1"),
				Port:    80,
				Network: net.Network_TCP,
			},
			Output: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x50},
		},
		{
			Input: net.Destination{
				Address: net.ParseAddress("2001:db8::1"),
				Port:    443,
				Network: net.Network_TCP,
			},
			Output: []byte{0x04, 0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0xbb},
		},
		{
			Input: net.Destination{
				Address: net.ParseAddress("akinokae.de"),
				Port:    443,
				Network: net.Network_TCP,
			},
			Output: []byte{0x03, 0x0b, 0x61, 0x6b, 0x69, 0x6e, 0x6f, 0x6b, 0x61, 0x65, 0x2e, 0x64, 0x65, 0x01, 0xbb},
		},
	}

	for _, testCase := range testCases {
		result := ParseDestination(testCase.Input)

		if r := cmp.Diff(result, testCase.Output); r != "" {
			t.Error(r)
		}
	}
}
