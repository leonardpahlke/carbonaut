package mockgeolocation

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
)

var PluginName plugin.Kind = "mockgeolocation"

type p struct{}

func New() p {
	return p{}
}

func (p) Get(cfg *staticenv.Config, resData *staticenv.InfraData) (*staticenv.Data, error) {
	return &staticenv.Data{
		Region:  "Frankfurt",
		Country: "Germany",
	}, nil
}
