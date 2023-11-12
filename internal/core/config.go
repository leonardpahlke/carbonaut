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
	"os"

	"carbonaut.cloud/internal/connector"
	"carbonaut.cloud/internal/metrics"
	"github.com/creasty/defaults"
	"github.com/gookit/validate"
	"gopkg.in/yaml.v3"
)

// Configuration is the Carbonaut configuration
type Configuration struct {
	// APIVersion string `yaml:"apiVersion" json:"apiVersion" validate:"required" default:"v1_alpha1"`
	Kind string `yaml:"kind" json:"kind" default:"carbonaut" validate:"required"`
	Meta Meta   `yaml:"meta" validate:"required" json:"meta"`
	Spec Spec   `yaml:"spec" validate:"required" json:"spec"`
}

type Meta struct {
	Name string `yaml:"name" json:"name" validate:"required" default:"carbonaut"`
}

type Spec struct {
	Provider []connector.Config `yaml:"provider" json:"provider"`
	Metrics  metrics.Config     `yaml:"server" json:"server" verify:"required"`
}

// Marshal transforms the referenced configuration struct to yaml string bytes
func (c *Configuration) Marshal() ([]byte, error) {
	return yaml.Marshal(c)
}

// UnmarshalConfiguration transforms the yaml string bytes to a configuration struct
func UnmarshalConfiguration(b []byte) (*Configuration, error) {
	c := Configuration{}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}
	return &c, nil
}

// ReadConfigurationFile reads the local Carbonaut configuration file
func ReadConfigurationFile(cfgFilePath string) (*Configuration, error) {
	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read configuration file: %w", err)
	}

	cfg := &Configuration{}
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := defaults.Set(cfg); err != nil {
		return nil, fmt.Errorf("could not set default values: %w", err)
	}

	slog.Info("configuration loaded", "path", cfgFilePath)
	return cfg, nil
}

// Validate validates the configuration
func (c *Configuration) Validate() error {
	v := validate.Struct(c)
	if !v.Validate() {
		return v.Errors
	}
	return nil
}
