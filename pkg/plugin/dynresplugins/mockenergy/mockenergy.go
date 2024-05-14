package mockenergy

import (
	"math/rand"

	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/types/dynres"
)

var PluginName plugin.Kind = "mockenergy"

type p struct{}

func New() p {
	return p{}
}

func (p) GetName() *plugin.Kind {
	return &PluginName
}

func (p) GetDynamicResourceData(cfg *dynres.Config, data *resource.StaticResData) (*resource.DynamicResData, error) {
	return &resource.DynamicResData{
		// Random CPU frequency between 1000 MHz and 3000 MHz
		CPUFrequency: 1000 + rand.Float64()*2000, // #nosec
		// Random energy consumption between 50 mW and 150 mW
		EnergyHostMilliwatt: rand.Intn(100) + 50, // #nosec
		// Random CPU load percentage from 0% to 100%
		CPULoadPercentage: rand.Float64() * 100, // #nosec
	}, nil
}
