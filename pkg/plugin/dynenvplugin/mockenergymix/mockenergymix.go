package mockenergymix

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/environment/dynenv"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
)

const PluginName plugin.Kind = "mockenergymix"

type p struct{}

func New() p {
	return p{}
}

func (p) Get(cfg dynenv.Config, data staticenv.Data) (dynenv.Data, error) {
	return dynenv.Data{
		SolarPercentage:        20.5,
		WindPercentage:         30.0,
		HydroPercentage:        15.0,
		NuclearPercentage:      10.0,
		FossilFuelsPercentage:  20.0,
		BioEnergyPercentage:    4.5,
		OtherSourcesPercentage: 0.0,
	}, nil
}
