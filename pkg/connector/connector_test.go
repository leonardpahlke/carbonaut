package connector

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"carbonaut.dev/pkg/plugin/dynenvplugins/mockenergymix"
	"carbonaut.dev/pkg/plugin/dynresplugins/mockenergy"
	"carbonaut.dev/pkg/plugin/staticresplugins/mockcloudplugin"
	"carbonaut.dev/pkg/provider"
	"carbonaut.dev/pkg/provider/data/account"
	"carbonaut.dev/pkg/provider/types/dynenv"
	"carbonaut.dev/pkg/provider/types/dynres"
	"carbonaut.dev/pkg/provider/types/staticres"
	"carbonaut.dev/pkg/util/logger"
)

func init() {
	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr, logger.DefaultOptions)))
}

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
					Plugin:       &mockcloudplugin.PluginName,
					AccessKeyEnv: &exampleAccessKeyA,
				},
				DynamicResConfig: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
			exampleAccountB: {
				StaticResConfig: &staticres.Config{
					Plugin:       &mockcloudplugin.PluginName,
					AccessKeyEnv: &exampleAccessKeyB,
				},
				DynamicResConfig: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
		},
		Environment: &provider.EnvConfig{
			DynamicEnvConfig: &dynenv.Config{
				Plugin: &mockenergymix.PluginName, AccessKeyEnv: &exampleAccessKeyC,
			},
		},
	}

	updatedProviderConfig = provider.Config{
		Resources: map[account.Name]provider.Res{
			exampleAccountA: {
				StaticResConfig: &staticres.Config{
					Plugin:       &mockcloudplugin.PluginName,
					AccessKeyEnv: &exampleAccessKeyB,
				},
				DynamicResConfig: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
			exampleAccountC: {
				StaticResConfig: &staticres.Config{
					Plugin:       &mockcloudplugin.PluginName,
					AccessKeyEnv: &exampleAccessKeyB,
				},
				DynamicResConfig: &dynres.Config{
					Plugin:    &mockenergy.PluginName,
					AccessKey: &exampleAccessKeyB,
				},
			},
		},
		Environment: &provider.EnvConfig{
			DynamicEnvConfig: &dynenv.Config{
				Plugin:       &mockenergymix.PluginName,
				AccessKeyEnv: &exampleAccessKeyA,
			},
		},
	}
)

func TestConnectorInit(t *testing.T) {
	if err := setMissingEnvVariables(&initialProviderConfig); err != nil {
		t.Error("internal test error", err)
	}

	if err := setMissingEnvVariables(&updatedProviderConfig); err != nil {
		t.Error("internal test error", err)
	}

	c, err := New(&Config{
		TimeoutSeconds: 10,
	}, &initialProviderConfig)
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
	if err := setMissingEnvVariables(&initialProviderConfig); err != nil {
		t.Error("internal test error", err)
	}

	connectorConfig := Config{
		TimeoutSeconds: 10,
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
			t.Fail()
		}
	}(t)
	go c.Run(stopChan, errChan)
	time.Sleep(2 * time.Second)
	stopChan <- 1
}

func TestConnectorCollect(t *testing.T) {
	if err := setMissingEnvVariables(&initialProviderConfig); err != nil {
		t.Error("internal test error", err)
	}

	connectorConfig := Config{
		TimeoutSeconds: 10,
	}

	c, err := New(&connectorConfig, &initialProviderConfig)
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

// PERFORMANCE TESTS

// BenchmarkCollect measures the performance of the Collect method.
func BenchmarkCollect(b *testing.B) {
	cfg := &Config{TimeoutSeconds: 10}

	if err := setMissingEnvVariables(&initialProviderConfig); err != nil {
		b.Error("internal test error", err)
	}

	c, err := New(cfg, &initialProviderConfig)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Collect()
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkLoadConfig measures the performance of the LoadConfig method.
func BenchmarkLoadConfig(b *testing.B) {
	cfg := &Config{TimeoutSeconds: 10}

	if err := setMissingEnvVariables(&initialProviderConfig); err != nil {
		b.Error("internal test error", err)
	}

	c, err := New(cfg, &initialProviderConfig)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			updatedConfig := updatedProviderConfig
			if err := c.LoadConfig(&updatedConfig); err != nil {
				b.Error("Error during LoadConfig: ", err)
			}
		}
	})
	b.StopTimer()
}

// mocked providers do not access these variables (use API keys etc.), but variables are checked if they are empty
func setMissingEnvVariables(cfg *provider.Config) error {
	if os.Getenv(*cfg.Environment.DynamicEnvConfig.AccessKeyEnv) == "" {
		if err := os.Setenv(*cfg.Environment.DynamicEnvConfig.AccessKeyEnv, "some value"); err != nil {
			return err
		}
	}

	for i := range cfg.Resources {
		if os.Getenv(*cfg.Resources[i].StaticResConfig.AccessKeyEnv) == "" {
			if err := os.Setenv(*cfg.Resources[i].StaticResConfig.AccessKeyEnv, "some value"); err != nil {
				return err
			}
		}
	}

	return nil
}
