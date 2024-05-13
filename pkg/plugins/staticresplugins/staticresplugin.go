package staticresplugins

import (
	"carbonaut.dev/pkg/plugins/staticresplugins/equinixplugin"
	"carbonaut.dev/pkg/plugins/staticresplugins/mockcloudplugin"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/types/staticres"
)

var plugins = map[plugin.Kind]staticres.Provider{
	mockcloudplugin.PluginName: mockcloudplugin.New(),
	equinixplugin.PluginName:   equinixplugin.New(),
}

func GetPlugin(identifier *plugin.Kind) (staticres.Provider, bool) {
	p, found := plugins[*identifier]
	return p, found
}
