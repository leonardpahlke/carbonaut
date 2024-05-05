package provider

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/environment"
	"carbonaut.dev/pkg/schema/provider/environment/dynenv"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
	"carbonaut.dev/pkg/schema/provider/resources"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

type Config struct {
	Resources   resources.Config   `json:"resources" yaml:"resources"`
	Environment environment.Config `json:"environment" yaml:"environment"`
}

// DATA
type Data map[plugin.AccountIdentifier][]AccountData

type AccountData struct {
	StaticResourceData     staticres.Data `json:"static_resource_data" yaml:"static_resource_data"`
	DynamicResourceData    dynres.Data    `json:"dynamic_resource_data" yaml:"dynamic_resource_data"`
	StaticEnvironmentData  staticenv.Data `json:"static_environment_data" yaml:"static_environment_data"`
	DynamicEnvironmentData dynenv.Data    `json:"dynamic_environment_data" yaml:"dynamic_environment_data"`
}
