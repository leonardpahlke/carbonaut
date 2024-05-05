package staticres

import (
	"carbonaut.dev/pkg/schema/plugin"
)

type Config struct {
	Plugin    plugin.Kind `json:"plugin" yaml:"plugin"`
	AccessKey string      `json:"access_key" yaml:"access_key"`
}

type Collector interface {
	ListResources(Config) ([]plugin.ResourceIdentifier, error)
	GetResource(Config, plugin.ResourceIdentifier) (Data, error)
}

// computer hardware data
type Data struct {
	Name      string `json:"name" yaml:"name"`
	IP        string `json:"ip" yaml:"ip"`
	CPUCores  int    `json:"cpu_cores" yaml:"cpu_cores"`
	MemoryMB  int    `json:"memory_mb" yaml:"memory_mb"`
	Arch      string `json:"arch" yaml:"arch"`
	StorageGB int    `json:"storage_gb" yaml:"storage_gb"`
}
