package staticres

import (
	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
)

type Config struct {
	Plugin       *plugin.Kind `json:"plugin"         yaml:"plugin"`
	AccessKeyEnv *string      `json:"access_key_env" yaml:"access_key_env"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetStaticResourceData(*project.Name, *resource.Name) (*resource.StaticResData, error)
	DiscoverStaticResourceIdentifiers(*project.Name) ([]*resource.Name, error)
	DiscoverProjectIdentifiers() ([]*project.Name, error)
}
