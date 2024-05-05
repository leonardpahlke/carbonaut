package mockenergy

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
	"carbonaut.dev/pkg/util/rnd"
)

const PluginName plugin.Kind = "mockenergy"

type p struct{}

func New() p {
	return p{}
}

func (p) Get(cfg dynres.Config, data staticres.Data) (dynres.Data, error) {
	return dynres.Data{
		// Random CPU frequency between 1000 MHz and 3000 MHz
		CPUFrequency: 1000 + rnd.RandFloat64()*2000,
		// Random energy consumption between 50 mW and 150 mW
		EnergyHostMilliwatt: rnd.RandIntn(100) + 50,
		// Random CPU load percentage from 0% to 100%
		CPULoadPercentage: rnd.RandFloat64() * 100,
	}, nil
}
