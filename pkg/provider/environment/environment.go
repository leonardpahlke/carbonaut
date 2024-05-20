package environment

type DynamicEnvData struct {
	SolarPercentage        float64 `json:"solar_percentage"         yaml:"solar_percentage"`
	WindPercentage         float64 `json:"wind_percentage"          yaml:"wind_percentage"`
	HydroPercentage        float64 `json:"hydro_percentage"         yaml:"hydro_percentage"`
	NuclearPercentage      float64 `json:"nuclear_percentage"       yaml:"nuclear_percentage"`
	GeothermalPercentage   float64 `json:"geothermal_percentage"    yaml:"geothermal_percentage"`
	GasPercentage          float64 `json:"gas_percentage"           yaml:"gas_percentage"`
	OilPercentage          float64 `json:"oil_percentage"           yaml:"oil_percentage"`
	BiomassPercentage      float64 `json:"biomass_percentage"       yaml:"biomass_percentage"`
	CoalPercentage         float64 `json:"coal_percentage"          yaml:"coal_percentage"`
	OtherSourcesPercentage float64 `json:"other_sources_percentage" yaml:"other_sources_percentage"`
	FossilFuelsPercentage  float64 `json:"fossil_fuels_percentage"  yaml:"fossil_fuels_percentage"`
	RenewablePercentage    float64 `json:"renewable_percentage"     yaml:"renewable_percentage"`
}
