package main

import (
	"log/slog"
	"os"

	"carbonaut.dev/pkg/config"
	"carbonaut.dev/pkg/connector"
	"carbonaut.dev/pkg/plugin/dynenvplugin/mockenergymix"
	"carbonaut.dev/pkg/plugin/dynresplugin/mockenergy"
	"carbonaut.dev/pkg/plugin/staticenvplugin/mockgeolocation"
	"carbonaut.dev/pkg/plugin/staticresplugin/mockcloudprovider"
	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/environment"
	"carbonaut.dev/pkg/schema/provider/environment/dynenv"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
	"carbonaut.dev/pkg/schema/provider/resources"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
	"carbonaut.dev/pkg/server"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

var (
	exampleAccessKeyA = "123"
	exampleAccessKeyB = "435"
	exampleAccessKeyC = "7654asdE2"
)

func main() {
	log := slog.Default()
	log.Info("Create a new Carbonaut Config")
	cfg := config.Config{
		Kind: "carbonaut",
		Meta: config.Meta{
			Name:     "carbonaut",
			LogLevel: "debug",
			Connector: &connector.Config{
				TimeoutSeconds: 10,
			},
		},
		Spec: config.Spec{
			Provider: &provider.Config{
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
			},
			Server: &server.Config{
				Port: 8088,
			},
		},
	}

	if err := defaults.Set(&cfg); err != nil {
		log.Error("failed to set defaults", "error", err)
		os.Exit(1)
	}

	y, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Error("failed to marshal config", "error", err)
		os.Exit(1)
	}

	log.Info("Write config to file")
	if err := os.WriteFile("config.yaml", y, 0o600); err != nil {
		log.Error("failed to write config to file", "error", err)
		os.Exit(1)
	}
	log.Info("Done, file written to config.yaml")
}
