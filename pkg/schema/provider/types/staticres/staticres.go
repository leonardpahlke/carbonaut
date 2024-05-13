package staticres

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/data/account/project"
	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
)

type Config struct {
	Plugin    *plugin.Kind `json:"plugin"     yaml:"plugin"`
	AccessKey *string      `json:"access_key" yaml:"access_key"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetStaticResourceData(*Config, *project.ID, *resource.ID) (*resource.StaticResData, error)
	DiscoverStaticResourceIdentifiers(*Config, *project.ID) ([]*resource.ID, error)
	DiscoverProjectIdentifiers(*Config) ([]*project.ID, error)
}
