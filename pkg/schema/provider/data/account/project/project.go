package project

import (
	"time"

	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
)

type (
	// defines the name of the project e.g. equinix-project-a
	ID   string
	Data map[resource.ID]*resource.Data
)

type Topology struct {
	Resources Resources `json:"resources"         yaml:"resources"`
	CreatedAt time.Time `json:"created_at"         yaml:"created_at"`
}

type Resources map[resource.ID]*resource.Topology
