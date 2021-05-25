package share

import (
	"fmt"
	"net"
)

func AvailablePort() int {
	// let the OS determine the available port
	l, err := net.Listen("tcp", "")
	if err != nil {
		// extremely rare case
		panic(fmt.Errorf("port unavailable: %w", err))
	}
	if err = l.Close(); err != nil {
		panic(err)
	}
	return l.Addr().(*net.TCPAddr).Port
}
