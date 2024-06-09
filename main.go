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

var (
	configFullPath string
	testRun        bool
)

const testRunSeconds = 5

func init() {
	flag.StringVar(&configFullPath, "c", "config.yaml", "Full path of the Carbonaut configuration file")
	flag.BoolVar(&testRun, "test-run", false, "If turned on Carbonaut runs a local test scenario and shut down. This can be used for verification purposes. Any provided configuration file is omitted!")
	flag.Parse()
}

func main() {
	exitChan := make(chan int)
	connectorErrChan := make(chan error)
	var cfg *config.Config
	if testRun {
		cfg = config.TestingConfiguration()
	} else {
		parsedCfg, err := config.ReadConfig(configFullPath)
		if err != nil {
			panic(fmt.Sprintf("could not read configuration file, err: %v", err))
		}
		cfg = parsedCfg
	}

	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr, &logger.Options{
		Level:       logger.GetLogLevel(cfg.Meta.LogLevel),
		TimeFormat:  time.DateTime,
		SrcFileMode: logger.ShortFile,
	})))

	if testRun {
		slog.Info("starting carbonaut in test mode, stopping carbonaut automatically", "stopping in", testRunSeconds)
		go func() {
			time.Sleep(testRunSeconds * time.Second)
			exitChan <- 1
		}()
		slog.Info("stopping carbonaut!")
	} else {
		slog.Info("starting carbonaut")
	}

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
	slog.Info("carbonaut shuts down, bye bye")
}
