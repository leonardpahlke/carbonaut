package iputils

import (
	"errors"
	"log/slog"
	"net"
	"os"
)

// GetOutboundIP tries to determine the outbound IP by making a UDP connection
// and returns the local network interface IP used for the connection.
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		slog.Error("error while getting outbound IP: %v", err)
		panic(err) // Depending on your application, you might prefer to return an error instead of panicking
	}
	defer conn.Close()

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		slog.Error("failed to type assert UDPAddr from LocalAddr")
		panic("failed to type assert UDPAddr from LocalAddr") // As above, consider how you handle this error
	}

	return localAddr.IP
}

func GetIPFromDomain(domain string) (string, error) {
	ip, err := net.LookupIP(os.Args[1])
	if err != nil {
		return "", err
	}
	if ip != nil {
		return ip[0].String(), nil
	}
	return "", errors.New("IP not found")
}
