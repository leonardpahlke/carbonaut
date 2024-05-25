package main

import (
	"log/slog"
	"os"
	"time"

	"carbonaut.dev/pkg/config"
	"carbonaut.dev/pkg/util/logger"
	"gopkg.in/yaml.v3"
)

func main() {
	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr, &logger.Options{
		Level:       slog.LevelInfo,
		TimeFormat:  time.DateTime,
		SrcFileMode: logger.ShortFile,
	})))
	slog.Info("Create a new Carbonaut Config")
	cfg := config.TestingConfiguration()

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
