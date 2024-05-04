package environment

import (
	"carbonaut.dev/pkg/connector/provider/environment/dynamic_env_collector"
	"carbonaut.dev/pkg/connector/provider/environment/static_env_collector"
)

type Data struct {
	Dynamic dynamic_env_collector.Data
	Static  static_env_collector.Data
}

type Config struct {
	DynamicEnvironment dynamic_env_collector.Config `json:"dynamic_environment"`
	StaticEnvironment  static_env_collector.Config  `json:"static_environment"`
}
