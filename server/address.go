package server

import (
	"fmt"
	"net"
	"strings"
)

func GetAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Get address : ", err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Error getting address")
			continue
		}

		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				ipad := v.IP.To4().String()
				if strings.Contains(ipad, "192") {
					return ipad
				}
			case *net.IPNet:
				ipad := v.IP.To4().String()
				if strings.Contains(ipad, "192") {
					return ipad
				}
			}

		}
	}
	return ""
}
