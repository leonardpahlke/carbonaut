/*
Copyright 2023 CARBONAUT AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"fmt"
	"log/slog"

	"carbonaut.cloud/internal/connector"
	"carbonaut.cloud/internal/metrics"
)

// Core defines the middleware of Carbonaut

type ConfigCore struct {
	// ConfigurationFilePath is the path to the configuration file
	ConfigurationFilePath string
}

type Core struct {
	cfg *Configuration
}

// New creates a new Carbonaut instance
func New(cfg *ConfigCore) (*Core, error) {
	slog.Info("read configuration...")
	configuration, err := ReadConfigurationFile(cfg.ConfigurationFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	slog.Info("validate configuration...")
	if err := configuration.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return &Core{
		cfg: configuration,
	}, nil
}

func (c *Core) Run() error {
	slog.Debug("starting Carbonaut")
	p, err := connector.CreateProviders(c.cfg.Spec.Provider, c.cfg.Spec.Metrics.CollectorName)
	if err != nil {
		return err
	}

	for i := range p {
		go p[i].ListenOnChannel()
	}

	slog.Info("starting metrics server...")
	metrics.Serve(c.cfg.Spec.Metrics, p)
	return nil
}
