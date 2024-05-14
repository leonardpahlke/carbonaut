package project

import (
	"time"

	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
)

type (
	// internal ID which get's counted up
	ID int32
	// defines the name of the project e.g. equinix-project-a
	Name string
	Data map[resource.Name]*resource.Data
)

// internal state
type Topology struct {
	Name              *Name     `json:"name"         yaml:"name"`
	Resources         Resources `json:"resources"         yaml:"resources"`
	CreatedAt         time.Time `json:"created_at"         yaml:"created_at"`
	ResourceIDCounter *int32    `json:"resource_id_counter"         yaml:"resource_id_counter"`
}

type Resources map[resource.ID]*resource.Topology

var NotFoundID ID = -1
