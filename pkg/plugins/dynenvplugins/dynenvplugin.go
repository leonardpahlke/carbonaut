package dynenvplugins

import (
	"carbonaut.dev/pkg/plugins/dynenvplugins/mockenergymix"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/types/dynenv"
)

var plugins = map[plugin.Kind]dynenv.Provider{
	mockenergymix.PluginName: mockenergymix.New(),
}

func GetPlugin(identifier *plugin.Kind) (dynenv.Provider, bool) {
	p, found := plugins[*identifier]
	return p, found
}
