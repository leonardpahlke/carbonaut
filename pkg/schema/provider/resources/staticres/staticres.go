package staticres

import (
	"carbonaut.dev/pkg/schema/plugin"
)

type Config struct {
	Plugin    plugin.Kind `json:"plugin"`
	AccessKey string      `json:"access_key"`
}

type Collector interface {
	ListResources(Config) ([]plugin.ResourceIdentifier, error)
	GetResource(Config, plugin.ResourceIdentifier) (Data, error)
}

// computer hardware data
type Data struct {
	Name      string `json:"name"`
	IP        string `json:"ip"`
	CPUCores  int    `json:"cpu_cores"`
	MemoryMB  int    `json:"memory_mb"`
	Arch      string `json:"arch"`
	StorageGB int    `json:"storage_gb"`
}
