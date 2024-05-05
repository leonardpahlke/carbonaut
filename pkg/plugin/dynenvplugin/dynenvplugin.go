package dynenvplugin

import (
	"fmt"

	"carbonaut.dev/pkg/plugin/dynenvplugin/mockenergymix"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/environment/dynenv"
)

var plugins = map[plugin.Kind]dynenv.Collector{
	mockenergymix.PluginName: mockenergymix.New(),
}

func GetPlugin(identifier plugin.Kind) (dynenv.Collector, error) {
	plugin, found := plugins[identifier]
	if found {
		return plugin, nil
	}
	return nil, fmt.Errorf("no plugin found with the name %s", identifier)
}

func GetPluginIdentifiers() []plugin.Kind {
	identifiers := make([]plugin.Kind, 0, len(plugins))
	for identifier := range plugins {
		identifiers = append(identifiers, identifier)
	}
	return identifiers
}
