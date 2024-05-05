package resources

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

// DATA
type Data map[plugin.AccountIdentifier][]AccountData

type AccountData struct {
	Static  staticres.Data
	Dynamic dynres.Data
}

// CONFIG
type Config map[plugin.AccountIdentifier]ResourceConfig

type ResourceConfig struct {
	StaticResource  staticres.Config `json:"static_resource"`
	DynamicResource dynres.Config    `json:"dynamic_resource"`
}
