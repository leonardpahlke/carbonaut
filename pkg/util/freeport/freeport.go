package freeport

import (
	"errors"
	"log"
	"net"
)

func Get() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func() {
		if closeErr := l.Close(); closeErr != nil {
			log.Printf("failed to close listener: %v", closeErr)
		}
	}()
	p, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.New("could not decode to net.TCPAddr")
	}
	return p.Port, nil
}
