package mockcloudprovider

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

var PluginName plugin.Kind = "mockcloudprovider"

type p struct{}

func New() p {
	return p{}
}

func (p) ListResources(cfg *staticres.Config) (*[]plugin.ResourceIdentifier, error) {
	return &[]plugin.ResourceIdentifier{"resource-a", "resource-b", "resource-c"}, nil
}

func (p) GetResource(cfg *staticres.Config, resource *plugin.ResourceIdentifier) (*staticres.Data, error) {
	return &staticres.Data{
		Name:      "machine-a",
		IP:        "0.0.0.0",
		CPUCores:  2,
		MemoryMB:  1000,
		Arch:      "arm64",
		StorageGB: 20,
	}, nil
}
