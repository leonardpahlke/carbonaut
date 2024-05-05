package environment

import (
	"carbonaut.dev/pkg/schema/provider/environment/dynenv"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
)

type Config struct {
	DynamicEnvironment dynenv.Config    `json:"dynamic_environment" yaml:"dynamic_environment"`
	StaticEnvironment  staticenv.Config `json:"static_environment" yaml:"static_environment"`
}
