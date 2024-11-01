package freeport_test

import (
	"log/slog"
	"net"
	"strconv"
	"testing"

	"carbonaut.dev/pkg/util/freeport"
)

func TestGetFreePort(t *testing.T) {
	port, err := freeport.Get()
	if err != nil {
		t.Error(err)
	}
	if port == 0 {
		t.Error("port == 0")
	}

	l, err := net.Listen("tcp", "localhost"+":"+strconv.Itoa(port))
	if err != nil {
		t.Error(err)
	}
	if closeErr := l.Close(); closeErr != nil {
		slog.Error("failed to close listener", "error", closeErr)
	}
}
