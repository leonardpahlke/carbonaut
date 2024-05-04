package dynamic_res_collector

import "carbonaut.dev/pkg/schema"

type Config struct {
	Plugin    schema.PluginName `json:"plugin"`
	AccessKey string            `json:"access_key"`
}

type Collector interface {
	Get(cfg Config) (Data, error)
}

// energy and utilization data
type Data struct {
	CPUFrequency        float32
	EnergyHostMilliwatt int
}
