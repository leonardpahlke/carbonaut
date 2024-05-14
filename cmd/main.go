package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"carbonaut.dev/pkg/config"
	"carbonaut.dev/pkg/connector"
	"carbonaut.dev/pkg/server"
)

var configFullPath string

func init() {
	flag.StringVar(&configFullPath, "c", "config.yaml", "Full path of the Carbonaut configuration file")
	flag.Parse()
}

func main() {
	exitChan := make(chan int)
	connectorErrChan := make(chan error)
	cfg, err := config.ReadConfig(configFullPath)
	if err != nil {
		panic(fmt.Sprintf("could not read configuration file, err: %v", err))
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: GetLogLevel(cfg.Meta.LogLevel),
	})
	log := slog.New(handler)
	log.Info("starting carbonaut", "config", cfg)

	c, err := connector.New(cfg.Meta.Connector, log, cfg.Spec.Provider)
	if err != nil {
		log.Error("could not initialize connector with provided configuration", "connector config", cfg.Meta.Connector, "provider config", cfg.Spec.Provider, "error", err)
		os.Exit(1)
	}

	log.Info("starting carbonaut server", "address", fmt.Sprintf("http://0.0.0.0:%d", cfg.Spec.Server.Port))
	s := server.New(c, log, exitChan)
	go s.Listen(cfg.Spec.Server)

	log.Info("starting carbonaut connector")
	c.Run(exitChan, connectorErrChan)
}

func GetLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
