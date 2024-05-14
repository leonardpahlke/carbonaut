package mockenergymix

import (
	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/data/environment"
	"carbonaut.dev/pkg/provider/types/dynenv"
)

var PluginName plugin.Kind = "mockenergymix"

type p struct{}

func New() p {
	return p{}
}

func (p) GetName() *plugin.Kind {
	return &PluginName
}

func (p) GetDynamicEnvironmentData(cfg *dynenv.Config, data *resource.Location) (*environment.DynamicEnvData, error) {
	return &environment.DynamicEnvData{
		SolarPercentage:        20.5,
		WindPercentage:         30.0,
		HydroPercentage:        15.0,
		NuclearPercentage:      10.0,
		FossilFuelsPercentage:  20.0,
		BioEnergyPercentage:    4.5,
		OtherSourcesPercentage: 0.0,
	}, nil
}
