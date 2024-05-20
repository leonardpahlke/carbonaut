package staticres

import (
	"carbonaut.dev/pkg/provider/plugin"
	"carbonaut.dev/pkg/provider/resource"
)

type Config struct {
	Plugin       *plugin.Kind `json:"plugin"         yaml:"plugin"`
	AccessKeyEnv *string      `json:"access_key_env" yaml:"access_key_env"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetStaticResourceData(*resource.ProjectName, *resource.ResourceName) (*resource.StaticResData, error)
	DiscoverStaticResourceIdentifiers(*resource.ProjectName) ([]*resource.ResourceName, error)
	DiscoverProjectIdentifiers() ([]*resource.ProjectName, error)
}
