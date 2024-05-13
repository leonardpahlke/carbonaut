package dynresplugins

import (
	"carbonaut.dev/pkg/plugins/dynresplugins/mockenergy"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/types/dynres"
)

var plugins = map[plugin.Kind]dynres.Provider{
	mockenergy.PluginName: mockenergy.New(),
}

func GetPlugin(identifier *plugin.Kind) (dynres.Provider, bool) {
	p, found := plugins[*identifier]
	return p, found
}
