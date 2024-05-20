package provider

import (
	"carbonaut.dev/pkg/provider/resource"
	"carbonaut.dev/pkg/provider/types/dynenv"
	"carbonaut.dev/pkg/provider/types/dynres"
	"carbonaut.dev/pkg/provider/types/staticres"
)

type Config struct {
	Resources   ResConfig  `json:"resources"   yaml:"resources"`
	Environment *EnvConfig `json:"environment" yaml:"environment"`
}

type ResConfig map[resource.AccountName]Res

type Res struct {
	StaticResConfig  *staticres.Config `json:"static_resource"  yaml:"static_resource"`
	DynamicResConfig *dynres.Config    `json:"dynamic_resource" yaml:"dynamic_resource"`
}

type EnvConfig struct {
	DynamicEnvConfig *dynenv.Config `json:"dynamic_environment" yaml:"dynamic_environment"`
}

type Data map[resource.AccountName]resource.AccountData
