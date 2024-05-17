package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"carbonaut.dev/pkg/config"
	"carbonaut.dev/pkg/connector"
	"carbonaut.dev/pkg/server"
	"carbonaut.dev/pkg/util/logger"
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

	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr, &logger.Options{
		Level:       logger.GetLogLevel(cfg.Meta.LogLevel),
		TimeFormat:  time.DateTime,
		SrcFileMode: logger.ShortFile,
	})))
	slog.Info("starting carbonaut", "config", cfg)

	c, err := connector.New(cfg.Meta.Connector, cfg.Spec.Provider)
	if err != nil {
		slog.Error("could not initialize connector with provided configuration", "connector config", cfg.Meta.Connector, "provider config", cfg.Spec.Provider, "error", err)
		os.Exit(1)
	}

	slog.Info("starting carbonaut server", "address", fmt.Sprintf("http://0.0.0.0:%d", cfg.Spec.Server.Port))
	s := server.New(c, exitChan)
	go s.Listen(cfg.Spec.Server)

	slog.Info("starting carbonaut connector")
	c.Run(exitChan, connectorErrChan)
}
