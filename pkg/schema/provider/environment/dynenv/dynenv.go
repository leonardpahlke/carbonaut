package dynenv

import (
	"carbonaut.dev/pkg/schema/plugin"
)

type Config struct {
	Plugin    plugin.Kind `json:"plugin"`
	AccessKey string      `json:"access_key"`
}

type Collector interface {
	Get(Config) (Data, error)
}

type Data struct {
	SolarPercentage        float64 `json:"solar_percentage"`
	WindPercentage         float64 `json:"wind_percentage"`
	HydroPercentage        float64 `json:"hydro_percentage"`
	NuclearPercentage      float64 `json:"nuclear_percentage"`
	FossilFuelsPercentage  float64 `json:"fossil_fuels_percentage"`
	BioEnergyPercentage    float64 `json:"bio_energy_percentage"`
	OtherSourcesPercentage float64 `json:"other_sources_percentage"`
}
