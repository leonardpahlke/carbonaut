package dynenv

import (
	"carbonaut.dev/pkg/provider/account/project/resource"
	"carbonaut.dev/pkg/provider/environment"
	"carbonaut.dev/pkg/provider/plugin"
)

type Config struct {
	Plugin       *plugin.Kind `json:"plugin"         yaml:"plugin"`
	AccessKeyEnv *string      `json:"access_key_env" yaml:"access_key_env"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetDynamicEnvironmentData(*resource.Location) (*environment.DynamicEnvData, error)
}
