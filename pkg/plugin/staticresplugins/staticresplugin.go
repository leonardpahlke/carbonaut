package staticresplugins

import (
	"fmt"

	"carbonaut.dev/pkg/plugin/staticresplugins/equinixplugin"
	"carbonaut.dev/pkg/plugin/staticresplugins/mockcloudplugin"
	"carbonaut.dev/pkg/provider/types/staticres"
)

func GetPlugin(cfg *staticres.Config) (staticres.Provider, error) {
	switch *cfg.Plugin {
	case mockcloudplugin.PluginName:
		return mockcloudplugin.New(cfg)
	case equinixplugin.PluginName:
		return equinixplugin.New(cfg)
	default:
		return nil, fmt.Errorf("plugin of kind %s not found", *cfg.Plugin)
	}
}
