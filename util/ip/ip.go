// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ip

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Define http headers.
const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	XClientIP     = "x-client-ip"
)

// GetLocalIP returns the non loopback local IP of the host.
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// RemoteIP returns the remote ip of the request.
func RemoteIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XClientIP); ip != "" {
		remoteAddr = ip
	} else if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

type IpFilter func(net.Interface, *net.IPNet) bool

// GetLocalIps 获取局域网ip地址
func GetLocalIps(filter IpFilter) []string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	var ips []string
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() == nil {
						continue
					}
					ipStr := ipnet.IP.String()

					if filter == nil {
						ips = append(ips, ipStr)
					} else if filter(netInterfaces[i], ipnet) {
						ips = append(ips, ipStr)
					}
				}
			}
		}

	}
	return ips
}

func GetLocalNoVmIps() []string {
	return GetLocalIps(func(netInterface net.Interface, ipNet *net.IPNet) bool {
		return !strings.Contains(netInterface.Name, "VMware")
	})
}
