package staticresplugin

import (
	"fmt"

	"carbonaut.dev/pkg/plugin/staticresplugin/mockcloudprovider"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

var plugins = map[plugin.Kind]staticres.Collector{
	mockcloudprovider.PluginName: mockcloudprovider.New(),
}

func GetPlugin(identifier *plugin.Kind) (staticres.Collector, error) {
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
