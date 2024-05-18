package dynenv

import (
	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/data/environment"
)

type Config struct {
	Plugin       *plugin.Kind `json:"plugin"         yaml:"plugin"`
	AccessKeyEnv *string      `json:"access_key_env" yaml:"access_key_env"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetDynamicEnvironmentData(*Config, *resource.Location) (*environment.DynamicEnvData, error)
}
