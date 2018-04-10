package utils

import (
	"fmt"
	"net"
	"strings"
)

//Utils to help get real ip like using this kind of format 'eth1:8088'
func GetRealAddrByNetwork(addr string) string {
	parts := strings.Split(addr, ":")

	if len(parts) != 2 {
		return addr
	}

	name := parts[0]
	port := parts[1]

	i, err := net.InterfaceByName(name)

	if err != nil {
		return addr
	}

	addrs, err := i.Addrs()

	if err != nil {
		return addr
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)

		if ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return fmt.Sprintf("%v:%v", ipnet.IP.String(), port)
			}
		}
	}

	return addr
}
