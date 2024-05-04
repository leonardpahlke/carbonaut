package resources

import (
	"carbonaut.dev/pkg/connector/provider/resources/dynamic_res_collector"
	"carbonaut.dev/pkg/connector/provider/resources/static_res_collector"
)

type AccountIdentifier string
type ResourceIdentifier string

// DATA
type Data map[AccountIdentifier][]AccountData

type AccountData struct {
	Static  static_res_collector.Data
	Dynamic dynamic_res_collector.Data
}

// CONFIG
type Config map[AccountIdentifier]ResourceConfig

type ResourceConfig struct {
	StaticResource  static_res_collector.Config  `json:"static_resource"`
	DynamicResource dynamic_res_collector.Config `json:"dynamic_resource"`
}
