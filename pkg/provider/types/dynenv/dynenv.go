package dynenv

import (
	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/data/environment"
)

type Config struct {
	Plugin    *plugin.Kind `json:"plugin"     yaml:"plugin"`
	AccessKey *string      `json:"access_key" yaml:"access_key"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetDynamicEnvironmentData(*Config, *resource.Location) (*environment.DynamicEnvData, error)
}
