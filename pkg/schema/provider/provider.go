package provider

import (
	"carbonaut.dev/pkg/schema/provider/data/account"
	"carbonaut.dev/pkg/schema/provider/types/dynenv"
	"carbonaut.dev/pkg/schema/provider/types/dynres"
	"carbonaut.dev/pkg/schema/provider/types/staticres"
)

type Config struct {
	Resources   ResConfig  `json:"resources"   yaml:"resources"`
	Environment *EnvConfig `json:"environment" yaml:"environment"`
}

type ResConfig map[account.ID]Res

type Res struct {
	StaticResConfig  *staticres.Config `json:"static_resource"  yaml:"static_resource"`
	DynamicResConfig *dynres.Config    `json:"dynamic_resource" yaml:"dynamic_resource"`
}

type EnvConfig struct {
	DynamicEnvConfig *dynenv.Config `json:"dynamic_environment" yaml:"dynamic_environment"`
}

type Topology map[account.ID]*account.Topology

type Data map[account.ID]account.Data
