/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 16:00
 */
package tools

import (
	"fmt"
	"net"
	"strings"
)

const (
	networkSplit = "@"
)

func ParseNetwork(str string) (network, addr string, err error) {
	if idx := strings.Index(str, networkSplit); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
		return
	} else {
		network = str[:idx]
		addr = str[idx+1:]
		return
	}
}

// GetContainerIP returns the container's actual IP address
// Used for service registration in multi-container deployments
func GetContainerIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no non-loopback IPv4 address found")
}

// GetServiceAddress returns the service address for etcd registration
// Replaces 0.0.0.0 with actual container IP if needed
func GetServiceAddress(network, addr string) string {
	if !strings.Contains(addr, "0.0.0.0") {
		return network + "@" + addr
	}

	containerIP, err := GetContainerIP()
	if err != nil {
		// Fallback to original address if IP detection fails
		return network + "@" + addr
	}

	actualAddr := strings.Replace(addr, "0.0.0.0", containerIP, 1)
	return network + "@" + actualAddr
}
