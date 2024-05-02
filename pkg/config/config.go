package config

import "fmt"

type Config struct {
	StaticResourceProviders    []StaticResourceConfig
	StaticEnvironmentProvider  PluginConfig
	DynamicEnvironmentProvider PluginConfig
}

type StaticResourceConfig struct {
	StaticResourceConfig  PluginConfig
	DynamicResourceConfig PluginConfig
}

type PluginConfig struct {
	PluginType PluginType
	AccessKey  string
}

type PluginType string

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
