package dynres

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

type Config struct {
	Plugin    plugin.Kind `json:"plugin" yaml:"plugin"`
	AccessKey string      `json:"access_key" yaml:"access_key"`
}

type Collector interface {
	Get(Config, staticres.Data) (Data, error)
}

// energy and utilization data
type Data struct {
	// in MHz
	CPUFrequency float64 `json:"cpu_frequency" yaml:"cpu_frequency"`
	// in milliwatts
	EnergyHostMilliwatt int `json:"energy_host_milliwatt" yaml:"energy_host_milliwatt"`
	// CPU load as a percentage
	CPULoadPercentage float64 `json:"cpu_load_percentage" yaml:"cpu_load_percentage"`
}
