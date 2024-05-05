package dynenv

import "carbonaut.dev/pkg/schema"

type Config struct {
	Plugin    schema.PluginName `json:"plugin"`
	AccessKey string            `json:"access_key"`
}

type Collector interface {
	Get(Config) (Data, error)
}

type Data struct{}
