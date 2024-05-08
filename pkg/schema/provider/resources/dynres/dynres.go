package dynres

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

type Config struct {
	Plugin    *plugin.Kind `json:"plugin"     yaml:"plugin"`
	AccessKey *string      `json:"access_key" yaml:"access_key"`
}

type Collector interface {
	Get(*Config, *staticres.Data) (*Data, error)
}

// energy and utilization data
type Data struct {
	CPUFrequency        float64 `json:"cpu_frequency"         yaml:"cpu_frequency"`
	EnergyHostMilliwatt int     `json:"energy_host_milliwatt" yaml:"energy_host_milliwatt"`
	CPULoadPercentage   float64 `json:"cpu_load_percentage"   yaml:"cpu_load_percentage"`
}
