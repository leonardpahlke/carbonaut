package staticenvplugin

import (
	"fmt"

	"carbonaut.dev/pkg/plugin/staticenvplugin/mockgeolocation"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
)

var plugins = map[plugin.Kind]staticenv.Collector{
	mockgeolocation.PluginName: mockgeolocation.New(),
}

func GetPlugin(identifier *plugin.Kind) (staticenv.Collector, error) {
	p, found := plugins[*identifier]
	if found {
		return p, nil
	}
	return nil, fmt.Errorf("no plugin found with the name %s", *identifier)
}

func GetPluginIdentifiers() []plugin.Kind {
	identifiers := make([]plugin.Kind, 0, len(plugins))
	for identifier := range plugins {
		identifiers = append(identifiers, identifier)
	}
	return identifiers
}
