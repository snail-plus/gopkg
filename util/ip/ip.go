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

// RemoteIP 返回请求的客户端 IP。
// 出于安全考虑，默认只使用 TCP 连接对端地址，不信任客户端可伪造的
// X-Forwarded-For / X-Real-IP / X-Client-IP 头。若应用部署在可信反向代理之后，
// 应使用 RemoteIPFromProxy 并确认代理会覆盖这些头。
func RemoteIP(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}
	return normalizeIP(host)
}

// RemoteIPFromProxy 在请求确实来自可信代理时使用：从代理注入的头中解析客户端 IP。
// X-Forwarded-For 取最左侧地址，并校验是否为合法 IP；无法解析时回退到 RemoteIP。
func RemoteIPFromProxy(req *http.Request) string {
	if v := req.Header.Get(XForwardedFor); v != "" {
		if idx := strings.IndexByte(v, ','); idx >= 0 {
			v = v[:idx]
		}
		v = strings.TrimSpace(v)
		if net.ParseIP(v) != nil {
			return normalizeIP(v)
		}
	}
	for _, name := range []string{XRealIP, XClientIP} {
		if v := strings.TrimSpace(req.Header.Get(name)); v != "" && net.ParseIP(v) != nil {
			return normalizeIP(v)
		}
	}
	return RemoteIP(req)
}

func normalizeIP(ip string) string {
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

type IpFilter func(net.Interface, *net.IPNet) bool

func getLocalAddrs(filter func(name string) bool) ([]net.IPNet, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var addresses []net.IPNet

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) == 0 {
			continue
		}
		if filter != nil && filter(netInterfaces[i].Name) {
			continue
		}
		addrs, _ := netInterfaces[i].Addrs()

		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() == nil {
					continue
				}
				addresses = append(addresses, *ipnet)
			}
		}
	}
	return addresses, nil
}

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

func isContainerInterface(name string) bool {
	n := strings.ToLower(strings.TrimSpace(name))
	return strings.HasPrefix(n, "docker") ||
		strings.HasPrefix(n, "br-") ||
		strings.HasPrefix(n, "veth")
}

func GetLocalNoDockerAddrs() ([]net.IPNet, error) {
	return getLocalAddrs(isContainerInterface)
}

func GetLocalNoDockerIps() ([]string, error) {
	addrs, _ := getLocalAddrs(isContainerInterface)
	var ips []string
	for _, addr := range addrs {
		ips = append(ips, addr.IP.String())
	}
	return ips, nil
}

func GetLocalNoVmIps() []string {
	return GetLocalIps(func(netInterface net.Interface, ipNet *net.IPNet) bool {
		return !strings.Contains(netInterface.Name, "VMware")
	})
}
