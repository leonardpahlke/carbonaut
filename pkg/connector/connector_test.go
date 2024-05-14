package connector

import (
	"log/slog"
	"testing"
	"time"

	"carbonaut.dev/pkg/plugins/dynenvplugins/mockenergymix"
	"carbonaut.dev/pkg/plugins/dynresplugins/mockenergy"
	"carbonaut.dev/pkg/plugins/staticresplugins/mockcloudplugin"
	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/data/account"
	"carbonaut.dev/pkg/schema/provider/types/dynenv"
	"carbonaut.dev/pkg/schema/provider/types/dynres"
	"carbonaut.dev/pkg/schema/provider/types/staticres"
)

// Adjust the initialProviderConfig and updatedProviderConfig to align with new Config structure
var (
	exampleAccessKeyA     = "123"
	exampleAccessKeyB     = "435"
	exampleAccessKeyC     = "7654asdE2"
	exampleAccountA       = account.Name("test-plugin-A")
	exampleAccountB       = account.Name("test-plugin-B")
	exampleAccountC       = account.Name("test-plugin-C")
	initialProviderConfig = provider.Config{
		Resources: map[account.Name]provider.Res{
			exampleAccountA: {
				StaticResConfig: &staticres.Config{
					Plugin:    &mockcloudplugin.PluginName,
					AccessKey: &exampleAccessKeyA,
				},
				DynamicResConfig: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
			exampleAccountB: {
				StaticResConfig: &staticres.Config{
					Plugin:    &mockcloudplugin.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
				DynamicResConfig: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
		},
		Environment: &provider.EnvConfig{
			DynamicEnvConfig: &dynenv.Config{
				Plugin: &mockenergymix.PluginName, AccessKey: &exampleAccessKeyC,
			},
		},
	}

	updatedProviderConfig = provider.Config{
		Resources: map[account.Name]provider.Res{
			exampleAccountA: {
				StaticResConfig: &staticres.Config{
					Plugin:    &mockcloudplugin.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
				DynamicResConfig: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
			exampleAccountC: {
				StaticResConfig: &staticres.Config{
					Plugin:    &mockcloudplugin.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
				DynamicResConfig: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
		},
		Environment: &provider.EnvConfig{
			DynamicEnvConfig: &dynenv.Config{
				Plugin:    &mockenergymix.PluginName,
				AccessKey: &exampleAccessKeyA,
			},
		},
	}
)

func TestConnectorInit(t *testing.T) {
	c, err := New(&Config{
		TimeoutSeconds: 10,
	}, slog.Default(), &initialProviderConfig)
	if err != nil {
		t.Error(err)
	}

	if err := c.LoadConfig(&updatedProviderConfig); err != nil {
		t.Error(err)
	}

	if len(c.state.T.Accounts) != len(updatedProviderConfig.Resources) {
		t.Error("state does not reflect configured resource accounts")
	}
	for aName := range updatedProviderConfig.Resources {
		found := false
		for i := range c.state.T.Accounts {
			if *c.state.T.Accounts[i].Name == aName {
				found = true
				if c.state.T.Accounts[i].Config == nil {
					t.Errorf("configuration of account not set")
				}
				continue
			}
		}
		if !found {
			t.Errorf("expected key %s was not found in map %v", aName, c.state.T)
		}
	}
}

func TestConnectorRun(t *testing.T) {
	connectorConfig := Config{
		TimeoutSeconds: 10,
	}

	c, err := New(&connectorConfig, slog.Default(), &initialProviderConfig)
	if err != nil {
		t.Error(err)
	}

	stopChan := make(chan int)
	errChan := make(chan error)
	go func(t *testing.T) {
		for e := range errChan {
			t.Error(e)
			t.Fail()
		}
	}(t)
	go c.Run(stopChan, errChan)
	time.Sleep(2 * time.Second)
	stopChan <- 1
}

func TestConnectorCollect(t *testing.T) {
	connectorConfig := Config{
		TimeoutSeconds: 10,
	}

	c, err := New(&connectorConfig, slog.Default(), &initialProviderConfig)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	stopChan := make(chan int)
	errChan := make(chan error)
	go func(t *testing.T) {
		for e := range errChan {
			t.Error(e)
		}
	}(t)
	go c.Run(stopChan, errChan)
	time.Sleep(1 * time.Second)
	d, err := c.Collect()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(d)
	stopChan <- 1
}

// BenchmarkCollect measures the performance of the Collect method.
func BenchmarkCollect(b *testing.B) {
	logger := slog.Default()
	config := &Config{TimeoutSeconds: 10}

	connector, err := New(config, logger, &initialProviderConfig)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := connector.Collect()
		if err != nil {
			b.Error(err)
		}
	}
}
