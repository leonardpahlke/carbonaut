package staticres

import (
	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
)

type Config struct {
	Plugin    *plugin.Kind `json:"plugin"     yaml:"plugin"`
	AccessKey *string      `json:"access_key" yaml:"access_key"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetStaticResourceData(*Config, *project.Name, *resource.Name) (*resource.StaticResData, error)
	DiscoverStaticResourceIdentifiers(*Config, *project.Name) ([]*resource.Name, error)
	DiscoverProjectIdentifiers(*Config) ([]*project.Name, error)
}
