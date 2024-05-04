package freeport

import (
	"errors"
	"net"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	p, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.New("could not decode to net.TCPAddr")
	}
	return p.Port, nil
}
