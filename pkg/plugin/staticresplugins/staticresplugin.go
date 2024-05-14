package staticresplugins

import (
	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/plugin/staticresplugins/equinixplugin"
	"carbonaut.dev/pkg/plugin/staticresplugins/mockcloudplugin"
	"carbonaut.dev/pkg/provider/types/staticres"
)

var plugins = map[plugin.Kind]staticres.Provider{
	mockcloudplugin.PluginName: mockcloudplugin.New(),
	equinixplugin.PluginName:   equinixplugin.New(),
}

func GetPlugin(identifier *plugin.Kind) (staticres.Provider, bool) {
	p, found := plugins[*identifier]
	return p, found
}
