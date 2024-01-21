package helpers

import (
	"fmt"
	"net"
)

func GetLocalIP() (string, error) {
	// Get the list of network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// Iterate through each network interface
	for _, iface := range interfaces {
		// Ignore loopback and non-up interfaces
		if iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagUp != 0 {
			// Get the list of IP addresses for the current interface
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			// Iterate through each IP address
			for _, addr := range addrs {
				// Check if the address is an IPv4 or IPv6 address
				switch v := addr.(type) {
				case *net.IPNet:
					// Exclude link-local and multicast addresses
					if !v.IP.IsLinkLocalUnicast() && !v.IP.IsMulticast() {
						return v.IP.String(), nil
					}
				case *net.IPAddr:
					return v.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("Unable to determine the local IP address")
}
