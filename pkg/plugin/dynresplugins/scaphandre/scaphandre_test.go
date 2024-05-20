package scaphandre_test

import (
	"fmt"
	"log/slog"
	"net"
	"testing"

	"carbonaut.dev/pkg/plugin/dynresplugins/scaphandre"
	"carbonaut.dev/pkg/provider/account/project/resource"
	"carbonaut.dev/pkg/provider/types/dynres"
)

func TestGetDynamicResourceData(t *testing.T) {
	port := 8080
	if !isPortInUse(port) {
		t.Log("there is nothing running on the specified port. skip test.")
		t.SkipNow()
	}
	p, err := scaphandre.New(&dynres.Config{
		Plugin:            scaphandre.PluginName,
		Endpoint:          fmt.Sprintf(":%d/metrics", port),
		AcceptInvalidCert: false,
	})
	if err != nil {
		t.Error(err)
	}
	_, err = p.GetDynamicResourceData(&resource.StaticResData{
		IPv4: "0.0.0.0",
	})
	if err != nil {
		t.Error(err)
	}
}

func isPortInUse(port int) bool {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return true
	}
	if closeErr := ln.Close(); closeErr != nil {
		slog.Error("failed to close listener: %v", closeErr)
	}
	return false
}
