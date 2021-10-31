package share

import (
	"fmt"
	"net"
	"strconv"
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

func PortValidate(s string) (int, error) {
	port, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}
	if port < 0 || port > 65536 {
		return -1, fmt.Errorf("invalid port number: %v", port)
	}
	return port, nil
}
