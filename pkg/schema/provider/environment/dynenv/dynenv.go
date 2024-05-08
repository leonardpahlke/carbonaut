package dynenv

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
)

type Config struct {
	Plugin    *plugin.Kind `json:"plugin"     yaml:"plugin"`
	AccessKey *string      `json:"access_key" yaml:"access_key"`
}

type Collector interface {
	Get(*Config, *staticenv.Data) (*Data, error)
}

type Data struct {
	SolarPercentage        float64 `json:"solar_percentage"         yaml:"solar_percentage"`
	WindPercentage         float64 `json:"wind_percentage"          yaml:"wind_percentage"`
	HydroPercentage        float64 `json:"hydro_percentage"         yaml:"hydro_percentage"`
	NuclearPercentage      float64 `json:"nuclear_percentage"       yaml:"nuclear_percentage"`
	FossilFuelsPercentage  float64 `json:"fossil_fuels_percentage"  yaml:"fossil_fuels_percentage"`
	BioEnergyPercentage    float64 `json:"bio_energy_percentage"    yaml:"bio_energy_percentage"`
	OtherSourcesPercentage float64 `json:"other_sources_percentage" yaml:"other_sources_percentage"`
}
