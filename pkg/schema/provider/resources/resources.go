package resources

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

// CONFIG
type Config map[plugin.AccountIdentifier]ResourceConfig

type ResourceConfig struct {
	StaticResource  *staticres.Config `json:"static_resource"  yaml:"static_resource"`
	DynamicResource *dynres.Config    `json:"dynamic_resource" yaml:"dynamic_resource"`
}
