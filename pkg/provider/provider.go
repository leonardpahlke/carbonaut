package provider

import (
	"carbonaut.dev/pkg/provider/data/account"
	"carbonaut.dev/pkg/provider/types/dynenv"
	"carbonaut.dev/pkg/provider/types/dynres"
	"carbonaut.dev/pkg/provider/types/staticres"
)

type Config struct {
	Resources   ResConfig  `json:"resources"   yaml:"resources"`
	Environment *EnvConfig `json:"environment" yaml:"environment"`
}

type ResConfig map[account.Name]Res

type Res struct {
	StaticResConfig  *staticres.Config `json:"static_resource"  yaml:"static_resource"`
	DynamicResConfig *dynres.Config    `json:"dynamic_resource" yaml:"dynamic_resource"`
}

type EnvConfig struct {
	DynamicEnvConfig *dynenv.Config `json:"dynamic_environment" yaml:"dynamic_environment"`
}

// internal state
type Topology struct {
	Accounts          map[account.ID]*account.Topology `json:"accounts"           yaml:"accounts"`
	AccountsIDCounter *int32                           `json:"project_id_counter" yaml:"project_id_counter"`
}

type Data map[account.Name]account.Data
