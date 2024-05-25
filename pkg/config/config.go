package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"carbonaut.dev/pkg/connector"
	"carbonaut.dev/pkg/plugin/dynenvplugins/mockenergymix"
	"carbonaut.dev/pkg/plugin/dynresplugins/mockenergy"
	"carbonaut.dev/pkg/plugin/staticresplugins/mockcloudplugin"
	"carbonaut.dev/pkg/provider"
	"carbonaut.dev/pkg/provider/resource"
	"carbonaut.dev/pkg/provider/types/dynenv"
	"carbonaut.dev/pkg/provider/types/dynres"
	"carbonaut.dev/pkg/provider/types/staticres"
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
	cleanedPath := filepath.Clean(os.ExpandEnv(path))
	absPath, err := filepath.Abs(cleanedPath)
	if err != nil {
		return nil, fmt.Errorf("could not resolve absolute path: %w", err)
	}

	matchYAML := regexp.MustCompile(`\.yaml$`)
	if !matchYAML.MatchString(absPath) {
		return nil, fmt.Errorf("invalid file type: %s", filepath.Ext(absPath))
	}

	// #nosec G304
	data, err := os.ReadFile(absPath)
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

func TestingConfiguration() *Config {
	exampleAccessKeyA := "123"
	exampleAccessKeyB := "435"
	exampleEndpoint := ":8080/metrics"
	cfg := Config{
		Kind: "carbonaut",
		Meta: Meta{
			Name:     "carbonaut",
			LogLevel: "debug",
			Connector: &connector.Config{
				TimeoutSeconds: 10,
			},
		},
		Spec: Spec{
			Provider: &provider.Config{
				Resources: map[resource.AccountName]provider.Res{
					"example-account-a": {
						StaticResConfig: &staticres.Config{
							Plugin:       &mockcloudplugin.PluginName,
							AccessKeyEnv: &exampleAccessKeyB,
						},
						DynamicResConfig: &dynres.Config{
							Plugin:   mockenergy.PluginName,
							Endpoint: exampleEndpoint,
						},
					},
					"example-account-c": {
						StaticResConfig: &staticres.Config{
							Plugin:       &mockcloudplugin.PluginName,
							AccessKeyEnv: &exampleAccessKeyB,
						},
						DynamicResConfig: &dynres.Config{
							Plugin:   mockenergy.PluginName,
							Endpoint: exampleAccessKeyB,
						},
					},
				},
				Environment: &provider.EnvConfig{
					DynamicEnvConfig: &dynenv.Config{
						Plugin:       &mockenergymix.PluginName,
						AccessKeyEnv: &exampleAccessKeyA,
					},
				},
			},
			Server: &server.Config{
				Port: 8088,
			},
		},
	}

	if err := defaults.Set(&cfg); err != nil {
		slog.Error("failed to set defaults", "error", err)
		os.Exit(1)
	}
	return &cfg
}
