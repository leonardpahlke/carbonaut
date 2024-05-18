package dynresplugins

import (
	"fmt"

	"carbonaut.dev/pkg/plugin/dynresplugins/mockenergy"
	"carbonaut.dev/pkg/provider/types/dynres"
)

func GetPlugin(cfg *dynres.Config) (dynres.Provider, error) {
	switch *cfg.Plugin {
	case mockenergy.PluginName:
		return mockenergy.New(cfg)
	default:
		return nil, fmt.Errorf("plugin of kind %s not found", *cfg.Plugin)
	}
}
