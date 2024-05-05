package connector

import (
	"log/slog"
	"testing"
	"time"

	"carbonaut.dev/pkg/plugin/dynenvplugin/mockenergymix"
	"carbonaut.dev/pkg/plugin/dynresplugin/mockenergy"
	"carbonaut.dev/pkg/plugin/staticenvplugin/mockgeolocation"
	"carbonaut.dev/pkg/plugin/staticresplugin/mockcloudprovider"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/environment"
	"carbonaut.dev/pkg/schema/provider/environment/dynenv"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
	"carbonaut.dev/pkg/schema/provider/resources"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

var (
	initialProviderConfig = provider.Config{
		Resources: map[plugin.AccountIdentifier]resources.ResourceConfig{
			"test-plugin-A": {
				StaticResource: staticres.Config{
					Plugin:    mockcloudprovider.PluginName,
					AccessKey: "321",
				},
				DynamicResource: dynres.Config{
					Plugin:    mockenergy.PluginName,
					AccessKey: "123",
				},
			},
			"test-plugin-B": {
				StaticResource: staticres.Config{
					Plugin:    mockcloudprovider.PluginName,
					AccessKey: "321",
				},
				DynamicResource: dynres.Config{
					Plugin:    mockenergy.PluginName,
					AccessKey: "321",
				},
			},
		},
		Environment: environment.Config{
			DynamicEnvironment: dynenv.Config{
				Plugin:    mockenergymix.PluginName,
				AccessKey: "321",
			},
			StaticEnvironment: staticenv.Config{
				Plugin:    mockgeolocation.PluginName,
				AccessKey: "123",
			},
		},
	}
	updatedProviderConfig = provider.Config{
		Resources: map[plugin.AccountIdentifier]resources.ResourceConfig{
			"test-plugin-A": {
				StaticResource: staticres.Config{
					Plugin:    mockcloudprovider.PluginName,
					AccessKey: "321",
				},
				DynamicResource: dynres.Config{
					Plugin:    "dynres-plug-A",
					AccessKey: "123",
				},
			},
			"test-plugin-C": {
				StaticResource: staticres.Config{
					Plugin:    "staticres-plug-C",
					AccessKey: "432",
				},
				DynamicResource: dynres.Config{
					Plugin:    mockenergy.PluginName,
					AccessKey: "456",
				},
			},
		},
		Environment: environment.Config{
			DynamicEnvironment: dynenv.Config{
				Plugin:    mockenergymix.PluginName,
				AccessKey: "321",
			},
			StaticEnvironment: staticenv.Config{
				Plugin:    "staticenv-plug",
				AccessKey: "123",
			},
		},
	}
)

func TestConnectorInit(t *testing.T) {
	connectorConfig := Config{
		TimeoutSeconds: 10,
		Log:            slog.Default(),
	}
	c, err := New(&connectorConfig, &initialProviderConfig)
	if err != nil {
		t.Error(err)
	}

	if err := c.LoadConfig(&updatedProviderConfig); err != nil {
		t.Error(err)
	}

	if len(c.state.Accounts) != len(updatedProviderConfig.Resources) {
		t.Error("state does not reflect configured resource accounts")
	}
	for accountID := range updatedProviderConfig.Resources {
		if _, exists := c.state.Accounts[accountID]; !exists {
			t.Errorf("expected key %s was not found in map %v", accountID, c.state.Accounts)
		}
	}
}

func TestConnectorRun(t *testing.T) {
	connectorConfig := Config{
		TimeoutSeconds: 10,
		Log:            slog.Default(),
	}

	c, err := New(&connectorConfig, &initialProviderConfig)
	if err != nil {
		t.Error(err)
	}

	stopChan := make(chan int)
	errChan := make(chan error)
	go func(t *testing.T) {
		for e := range errChan {
			t.Error(e)
		}
	}(t)
	go c.Run(stopChan, errChan)
	time.Sleep(2 * time.Second)
	stopChan <- 1
}
