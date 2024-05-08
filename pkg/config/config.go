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
	Kind string `default:"carbonaut" json:"kind"         validate:"required" yaml:"kind"`
	Meta Meta   `json:"meta"         validate:"required" yaml:"meta"`
	Spec Spec   `json:"spec"         validate:"required" yaml:"spec"`
}

type Meta struct {
	Name      string            `default:"carbonaut" json:"name"      validate:"required" yaml:"name"`
	LogLevel  string            `default:"info"      json:"log_level" yaml:"log_level"`
	Connector *connector.Config `json:"connector"    yaml:"connector"`
}

type Spec struct {
	Provider *provider.Config `json:"provider" yaml:"provider"`
	Server   *server.Config   `json:"server"   verify:"required" yaml:"server"`
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
