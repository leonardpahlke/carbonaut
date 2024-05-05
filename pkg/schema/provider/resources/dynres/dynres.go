package dynres

import (
	"carbonaut.dev/pkg/schema/plugin"
)

type Config struct {
	Plugin    plugin.Kind `json:"plugin"`
	AccessKey string      `json:"access_key"`
}

type Collector interface {
	Get(cfg Config) (Data, error)
}

// energy and utilization data
type Data struct {
	// in MHz
	CPUFrequency float64 `json:"cpu_frequency"`
	// in milliwatts
	EnergyHostMilliwatt int `json:"energy_host_milliwatt"`
	// CPU load as a percentage
	CPULoadPercentage float64 `json:"cpu_load_percentage"`
}
