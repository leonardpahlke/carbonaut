package main

import (
	"log/slog"
	"os"
	"time"

	"carbonaut.dev/pkg/config"
	"carbonaut.dev/pkg/connector"
	"carbonaut.dev/pkg/plugin/dynenvplugins/mockenergymix"
	"carbonaut.dev/pkg/plugin/dynresplugins/mockenergy"
	"carbonaut.dev/pkg/plugin/staticresplugins/mockcloudplugin"
	"carbonaut.dev/pkg/provider"
	"carbonaut.dev/pkg/provider/data/account"
	"carbonaut.dev/pkg/provider/types/dynenv"
	"carbonaut.dev/pkg/provider/types/dynres"
	"carbonaut.dev/pkg/provider/types/staticres"
	"carbonaut.dev/pkg/server"
	"carbonaut.dev/pkg/util/logger"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

var (
	exampleAccessKeyA = "123"
	exampleAccessKeyB = "435"
)

func main() {
	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr, &logger.Options{
		Level:       slog.LevelInfo,
		TimeFormat:  time.DateTime,
		SrcFileMode: logger.ShortFile,
	})))
	slog.Info("Create a new Carbonaut Config")
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
				Resources: map[account.Name]provider.Res{
					"example-account-a": {
						StaticResConfig: &staticres.Config{
							Plugin:       &mockcloudplugin.PluginName,
							AccessKeyEnv: &exampleAccessKeyB,
						},
						DynamicResConfig: &dynres.Config{
							Plugin:    &mockenergy.PluginName,
							AccessKey: &exampleAccessKeyB,
						},
					},
					"example-account-c": {
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
			},
			Server: &server.Config{
				Port: 8088,
			},
		},
	}

	if err := defaults.Set(&cfg); err != nil {
		slog.Error("failed to set defaults", "error", err)
		os.Exit(1)
	}

	y, err := yaml.Marshal(&cfg)
	if err != nil {
		slog.Error("failed to marshal config", "error", err)
		os.Exit(1)
	}

	slog.Info("Write config to file")
	if err := os.WriteFile("config.yaml", y, 0o600); err != nil {
		slog.Error("failed to write config to file", "error", err)
		os.Exit(1)
	}
	slog.Info("Done, file written to config.yaml")
}
