package config

import (
	"fmt"

	"carbonaut.dev/pkg/connector"
	"carbonaut.dev/pkg/connector/provider"
)

type Config struct {
	Meta Meta             `json:"meta"`
	Spec *provider.Config `json:"spec"`
}

type Meta struct {
	Connector *connector.Config `json:"connector"`
}

func ReadConfig(path string) (*Config, error) {
	fmt.Printf("collect config from path: %s\n", path)

	// TODO: read config from local file

	if err := ValidateConfig(); err != nil {
		return nil, err
	}
	return &Config{}, nil
}

// TODO: validate configuration
// TODO: this may be limited to X number of resources. This depends on load testing results.
func ValidateConfig() error {
	return nil
}
