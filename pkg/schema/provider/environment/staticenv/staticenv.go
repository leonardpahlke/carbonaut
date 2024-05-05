package staticenv

import (
	"carbonaut.dev/pkg/schema/plugin"
)

type Config struct {
	Plugin    plugin.Kind `json:"plugin" yaml:"plugin"`
	AccessKey string      `json:"access_key" yaml:"access_key"`
}

type InfraData struct {
	IP string
}

type Collector interface {
	Get(Config, InfraData) (Data, error)
}

// location data
type Data struct {
	Region  string `json:"region" yaml:"region"`
	Country string `json:"country" yaml:"country"`
}
