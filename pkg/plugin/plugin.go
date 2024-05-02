package plugin

import "carbonaut.dev/pkg/data"

type ObserveConfig struct {
	// Defines how long to wait between checking resources
	ObserveTimeout int
}

type StaticResourcePlugin interface {
	Observe(ObserveConfig) data.StaticResourceData
}

type DynamicResourcePlugin interface {
	Observe() data.DynamicResourceData
}

type StaticEnvironmentPlugin interface {
	Observe() data.StaticEnvironmentData
}

type DynamicEnvironmentPlugin interface {
	Observe() data.DynamicEnvironmentData
}
