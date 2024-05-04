package static_res_collector

import "carbonaut.dev/pkg/schema"

type Config struct {
	Plugin    schema.PluginName `json:"plugin"`
	AccessKey string            `json:"access_key"`
}

type Collector interface {
	ListResources(Config) (schema.Resource, error)
	GetResource(Config, schema.Resource) (Data, error)
}

// computer hardware data
type Data struct {
	IP        string
	CPUCores  int
	MemoryMB  int
	Arch      string
	StorageGB int
}
