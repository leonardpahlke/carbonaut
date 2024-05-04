package resources

import (
	"carbonaut.dev/pkg/connector/provider/resources/dynres"
	"carbonaut.dev/pkg/connector/provider/resources/staticres"
)

type (
	AccountIdentifier  string
	ResourceIdentifier string
)

// DATA
type Data map[AccountIdentifier][]AccountData

type AccountData struct {
	Static  staticres.Data
	Dynamic dynres.Data
}

// CONFIG
type Config map[AccountIdentifier]ResourceConfig

type ResourceConfig struct {
	StaticResource  staticres.Config `json:"static_resource"`
	DynamicResource dynres.Config    `json:"dynamic_resource"`
}
