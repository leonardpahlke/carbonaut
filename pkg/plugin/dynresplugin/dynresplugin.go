package dynresplugin

import (
	"fmt"

	"carbonaut.dev/pkg/plugin/dynresplugin/mockenergy"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
)

var plugins = map[plugin.Kind]dynres.Collector{
	mockenergy.PluginName: mockenergy.New(),
}

func GetPlugin(identifier plugin.Kind) (dynres.Collector, error) {
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
