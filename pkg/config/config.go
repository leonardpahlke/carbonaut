package config

import (
	"fmt"
	"os"

	"carbonaut.dev/pkg/connector"
	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/server"
	"github.com/creasty/defaults"
	"github.com/gookit/validate"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Kind string `yaml:"kind" json:"kind" default:"carbonaut" validate:"required"`
	Meta Meta   `yaml:"meta" validate:"required" json:"meta"`
	Spec Spec   `yaml:"spec" validate:"required" json:"spec"`
}

type Meta struct {
	Name      string            `yaml:"name" json:"name" validate:"required" default:"carbonaut"`
	LogLevel  string            `yaml:"log_level" json:"log_level" default:"info"`
	Connector *connector.Config `yaml:"connector" json:"connector"`
}

type Spec struct {
	Provider *provider.Config `yaml:"provider" json:"provider"`
	Server   *server.Config   `yaml:"server" json:"server" verify:"required"`
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read Config file: %w", err)
	}

	cfg := &Config{}
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := defaults.Set(cfg); err != nil {
		return nil, fmt.Errorf("could not set default values: %w", err)
	}
	return cfg, nil
}

// TODO: validate config
// TODO: this may be limited to X number of resources. This depends on load testing results.
func ValidateConfig() error {
	return nil
}

// Marshal transforms the referenced config struct to yaml string bytes
func (c *Config) Marshal() ([]byte, error) {
	return yaml.Marshal(c)
}

// UnmarshalConfig transforms the yaml string bytes to a Config struct
func UnmarshalConfig(b []byte) (*Config, error) {
	c := Config{}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}
	return &c, nil
}

// Validate validates the Config
func (c *Config) Validate() error {
	v := validate.Struct(c)
	if !v.Validate() {
		return v.Errors
	}
	return nil
}
