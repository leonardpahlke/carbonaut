package environment

import (
	"carbonaut.dev/pkg/connector/provider/environment/dynenv"
	"carbonaut.dev/pkg/connector/provider/environment/staticenv"
)

type Data struct {
	Dynamic dynenv.Data
	Static  staticenv.Data
}

type Config struct {
	DynamicEnvironment dynenv.Config    `json:"dynamic_environment"`
	StaticEnvironment  staticenv.Config `json:"static_environment"`
}
