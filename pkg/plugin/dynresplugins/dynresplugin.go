package dynresplugins

import (
	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/plugin/dynresplugins/mockenergy"
	"carbonaut.dev/pkg/provider/types/dynres"
)

var plugins = map[plugin.Kind]dynres.Provider{
	mockenergy.PluginName: mockenergy.New(),
}

func GetPlugin(identifier *plugin.Kind) (dynres.Provider, bool) {
	p, found := plugins[*identifier]
	return p, found
}
