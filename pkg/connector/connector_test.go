package connector

import (
	"log/slog"
	"testing"

	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/environment"
	"carbonaut.dev/pkg/schema/provider/environment/dynenv"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
	"carbonaut.dev/pkg/schema/provider/resources"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

func TestConnectorInit(t *testing.T) {
	connectorConfig := Config{
		TimeoutSeconds: 10,
		Log:            slog.Default(),
	}
	initialProviderConfig := provider.Config{
		Resources: map[resources.AccountIdentifier]resources.ResourceConfig{
			"test-plugin-A": {
				StaticResource: staticres.Config{
					Plugin:    "staticres-plug-A",
					AccessKey: "321",
				},
				DynamicResource: dynres.Config{
					Plugin:    "dynres-plug-A",
					AccessKey: "123",
				},
			},
			"test-plugin-B": {
				StaticResource: staticres.Config{
					Plugin:    "staticres-plug-B",
					AccessKey: "321",
				},
				DynamicResource: dynres.Config{
					Plugin:    "dynres-plug-B",
					AccessKey: "321",
				},
			},
		},
		Environment: environment.Config{
			DynamicEnvironment: dynenv.Config{
				Plugin:    "dynenv-plug",
				AccessKey: "321",
			},
			StaticEnvironment: staticenv.Config{
				Plugin:    "staticenv-plug",
				AccessKey: "123",
			},
		},
	}
	c, err := New(&connectorConfig, &initialProviderConfig)
	if err != nil {
		t.Error(err)
	}

	updatedProviderConfig := provider.Config{
		Resources: map[resources.AccountIdentifier]resources.ResourceConfig{
			"test-plugin-A": {
				StaticResource: staticres.Config{
					Plugin:    "staticres-plug-A",
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
					Plugin:    "dynres-plug-C",
					AccessKey: "456",
				},
			},
		},
		Environment: environment.Config{
			DynamicEnvironment: dynenv.Config{
				Plugin:    "dynenv-plug",
				AccessKey: "321",
			},
			StaticEnvironment: staticenv.Config{
				Plugin:    "staticenv-plug",
				AccessKey: "123",
			},
		},
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
