package scaphandre_test

import (
	"testing"

	"carbonaut.dev/pkg/plugin/dynresplugins/scaphandre"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/types/dynres"
)

func TestGetDynamicResourceData(t *testing.T) {
	p, err := scaphandre.New(&dynres.Config{
		Plugin:            scaphandre.PluginName,
		Endpoint:          ":8080/metrics",
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
