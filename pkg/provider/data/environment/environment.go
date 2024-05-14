package environment

type DynamicEnvData struct {
	SolarPercentage        float64 `json:"solar_percentage"         yaml:"solar_percentage"`
	WindPercentage         float64 `json:"wind_percentage"          yaml:"wind_percentage"`
	HydroPercentage        float64 `json:"hydro_percentage"         yaml:"hydro_percentage"`
	NuclearPercentage      float64 `json:"nuclear_percentage"       yaml:"nuclear_percentage"`
	FossilFuelsPercentage  float64 `json:"fossil_fuels_percentage"  yaml:"fossil_fuels_percentage"`
	BioEnergyPercentage    float64 `json:"bio_energy_percentage"    yaml:"bio_energy_percentage"`
	OtherSourcesPercentage float64 `json:"other_sources_percentage" yaml:"other_sources_percentage"`
}
