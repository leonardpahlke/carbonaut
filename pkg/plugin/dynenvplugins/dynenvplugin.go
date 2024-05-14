package dynenvplugins

import (
	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/plugin/dynenvplugins/mockenergymix"
	"carbonaut.dev/pkg/provider/types/dynenv"
)

var plugins = map[plugin.Kind]dynenv.Provider{
	mockenergymix.PluginName: mockenergymix.New(),
}

func GetPlugin(identifier *plugin.Kind) (dynenv.Provider, bool) {
	p, found := plugins[*identifier]
	return p, found
}
