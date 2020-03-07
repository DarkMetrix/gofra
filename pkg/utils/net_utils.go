package utils

import (
	"fmt"
	"net"
)

//Utils to help get real ip like using this kind of format 'eth1:8088'
func GetRealAddrByNetwork(addr string) string {
	host, port, err := net.SplitHostPort(addr)

	if err != nil {
		return addr
	}

	ip, err := GetRealIPv4ByNetwork(host)

	if err != nil {
		return addr
	} else {
		return fmt.Sprintf("%v:%v", ip, port)
	}
}

//Utils to help get real ip using name like 'eth1'
func GetRealIPv4ByNetwork(name string) (string, error) {
	i, err := net.InterfaceByName(name)

	if err != nil {
		return "", err
	}

	addrs, err := i.Addrs()

	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)

		if ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "0.0.0.0", nil
}
