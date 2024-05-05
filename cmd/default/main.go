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
	"carbonaut.dev/pkg/schema/plugin"
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
				Resources: map[plugin.AccountIdentifier]resources.ResourceConfig{
					"a-cloudprovider-account": {
						StaticResource: staticres.Config{
							Plugin:    mockcloudprovider.PluginName,
							AccessKey: "32-some-kind-of-key-23",
						},
						DynamicResource: dynres.Config{
							Plugin:    mockenergy.PluginName,
							AccessKey: "12-some-kind-of-key-21",
						},
					},
					"b-cloudprovider-account": {
						StaticResource: staticres.Config{
							Plugin:    mockcloudprovider.PluginName,
							AccessKey: "456-some-kind-of-key-452",
						},
						DynamicResource: dynres.Config{
							Plugin:    mockenergy.PluginName,
							AccessKey: "2-some-kind-of-key-1",
						},
					},
				},
				Environment: environment.Config{
					DynamicEnvironment: dynenv.Config{
						Plugin:    mockenergymix.PluginName,
						AccessKey: "some-kind-of-key-22",
					},
					StaticEnvironment: staticenv.Config{
						Plugin:    mockgeolocation.PluginName,
						AccessKey: "22-some-kind-of-key",
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
