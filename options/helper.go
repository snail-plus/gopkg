// Copyright 2024 eve.  All rights reserved.

package options

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Define unit constant.
const (
	_   = iota // ignore onex.iota
	KiB = 1 << (10 * iota)
	MiB
	GiB
	TiB
)

func join(prefixs ...string) string {
	joined := strings.Join(prefixs, ".")
	if joined != "" {
		joined += "."
	}

	return joined
}

// ValidateAddress takes an address as a string and validates it.
// If the input address is not in a valid :port or IP:port format, it returns an error.
// It also checks if the host part of the address is a valid IP address and if the port number is valid.
func ValidateAddress(addr string) error {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("%q is not in a valid format (:port or ip:port): %w", addr, err)
	}
	if _, err := ParsePort(port, true); err != nil {
		return fmt.Errorf("%q is not a valid number", port)
	}

	return nil
}

// CreateListener create net listener by given address and returns it and port.
func CreateListener(addr string) (net.Listener, int, error) {
	network := "tcp"

	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to listen on %v: %w", addr, err)
	}

	// get port
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		_ = ln.Close()

		return nil, 0, fmt.Errorf("invalid listen address: %q", ln.Addr().String())
	}

	return ln, tcpAddr.Port, nil
}

func ParsePort(port string, allowZero bool) (int, error) {
	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, err
	}
	if portInt == 0 && !allowZero {
		return 0, errors.New("0 is not a valid port number")
	}
	return int(portInt), nil
}
