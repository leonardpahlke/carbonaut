package provider

import (
	"carbonaut.dev/pkg/connector/provider/environment"
	"carbonaut.dev/pkg/connector/provider/resources"
)

type Config struct {
	Resources   resources.Config   `json:"resources"`
	Environment environment.Config `json:"environment"`
}

type Data struct {
	Resources   resources.Data   `json:"resources"`
	Environment environment.Data `json:"environment"`
}
