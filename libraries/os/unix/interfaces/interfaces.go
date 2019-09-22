package interfaces

import (
	"fmt"
	"net"
	"strings"

	"../data"
)

type Interface struct {
	Name string
	IP   string
	Up   bool
}

const uninitialized = "Unknown"

var interfaces []Interface

func Header() *data.Header {
	return &header
}

func Data() []Interface {
	channel := make(chan *Interface)

	for i := 0; i < len(interfaces); i++ {
		go call(&interfaces[i], channel)
	}

	// ensure we get all responses back
	for range interfaces {
		_ = <-channel
	}

	return interfaces
}

func call(iface *Interface, channel chan *Interface) {
	addresses, up := getInterfaceAddresses(iface.Name)

	iface.IP = strings.Join(addresses, ",")
	iface.Up = up

	channel <- iface
}

func getInterfaceAddresses(interfaceName string) ([]string, bool) {
	addresses := make([]string, 0)

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return []string{"N/A"}, false
	}

	addrs, err := iface.Addrs()

	if err != nil {
		return []string{"N/A"}, false
	}

	// handle err
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		addresses = append(addresses, fmt.Sprintf("%s", ip))
	}

	var up bool
	if iface.Flags&net.FlagUp != 0 {
		up = true
	} else {
		up = false
	}

	return addresses, up
}
