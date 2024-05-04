package iputils

import (
	"errors"
	"net"
	"os"

	"log/slog"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		slog.Error("error while getting outbound ip: %v", err)
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GetIPFromDomain(domain string) (string, error) {
	ip, _ := net.LookupIP(os.Args[1])
	if ip != nil {
		return ip[0].String(), nil
	}
	return "", errors.New("IP not found")
}
