package static_env_collector

import "carbonaut.dev/pkg/schema"

type Config struct {
	Plugin    schema.PluginName `json:"plugin"`
	AccessKey string            `json:"access_key"`
}

type Collector interface {
	Get(Config) (Data, error)
}

// location data
type Data struct {
	Region  string
	Country string
}
