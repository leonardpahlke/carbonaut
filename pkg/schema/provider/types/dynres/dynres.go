package dynres

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
)

type Config struct {
	Plugin    *plugin.Kind `json:"plugin"     yaml:"plugin"`
	AccessKey *string      `json:"access_key" yaml:"access_key"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetDynamicResourceData(*Config, *resource.StaticResData) (*resource.DynamicResData, error)
}
