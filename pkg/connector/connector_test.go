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

// Adjust the initialProviderConfig and updatedProviderConfig to align with new Config structure
var (
	exampleAccessKeyA     = "123"
	exampleAccessKeyB     = "435"
	exampleAccessKeyC     = "7654asdE2"
	examplePluginA        = plugin.Kind("dynres-plug-A")
	examplePluginB        = plugin.Kind("dynres-plug-B")
	examplePluginC        = plugin.Kind("dynres-plug-C")
	initialProviderConfig = provider.Config{
		Resources: &resources.Config{
			"test-plugin-A": resources.ResourceConfig{
				StaticResource: &staticres.Config{
					Plugin:    &mockcloudprovider.PluginName,
					AccessKey: &exampleAccessKeyA,
				},
				DynamicResource: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
			"test-plugin-B": resources.ResourceConfig{
				StaticResource: &staticres.Config{
					Plugin:    &mockcloudprovider.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
				DynamicResource: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
		},
		Environment: &environment.Config{
			DynamicEnvironment: &dynenv.Config{
				Plugin:    &mockenergymix.PluginName,
				AccessKey: &exampleAccessKeyC,
			},
			StaticEnvironment: &staticenv.Config{
				Plugin:    &mockgeolocation.PluginName,
				AccessKey: &exampleAccessKeyA,
			},
		},
	}
	updatedProviderConfig = provider.Config{
		Resources: &resources.Config{
			"test-plugin-A": resources.ResourceConfig{
				StaticResource: &staticres.Config{
					Plugin:    &mockcloudprovider.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
				DynamicResource: &dynres.Config{
					Plugin:    &examplePluginA,
					AccessKey: &exampleAccessKeyB,
				},
			},
			"test-plugin-C": resources.ResourceConfig{
				StaticResource: &staticres.Config{
					Plugin:    &examplePluginC,
					AccessKey: &exampleAccessKeyC,
				},
				DynamicResource: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyC,
				},
			},
		},
		Environment: &environment.Config{
			DynamicEnvironment: &dynenv.Config{
				Plugin:    &mockenergymix.PluginName,
				AccessKey: &exampleAccessKeyA,
			},
			StaticEnvironment: &staticenv.Config{
				Plugin:    &examplePluginB,
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

	if len(c.state.Accounts) != len(*updatedProviderConfig.Resources) {
		t.Error("state does not reflect configured resource accounts")
	}
	for accountID := range *updatedProviderConfig.Resources {
		if _, exists := c.state.Accounts[accountID]; !exists {
			t.Errorf("expected key %s was not found in map %v", accountID, c.state.Accounts)
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
