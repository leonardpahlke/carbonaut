package dynenvplugins

import (
	"fmt"

	"carbonaut.dev/pkg/plugin/dynenvplugins/electricitymaps"
	"carbonaut.dev/pkg/plugin/dynenvplugins/mockenergymix"
	"carbonaut.dev/pkg/provider/types/dynenv"
)

func GetPlugin(cfg *dynenv.Config) (dynenv.Provider, error) {
	switch *cfg.Plugin {
	case mockenergymix.PluginName:
		return mockenergymix.New(cfg)
	case electricitymaps.PluginName:
		return electricitymaps.New(cfg)
	default:
		return nil, fmt.Errorf("plugin of kind %s not found", *cfg.Plugin)
	}
}
