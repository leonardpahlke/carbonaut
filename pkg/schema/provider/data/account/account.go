package account

import (
	"time"

	"carbonaut.dev/pkg/schema/provider/data/account/project"
	"carbonaut.dev/pkg/schema/provider/types/staticres"
)

type (
	// internal ID which get's counted up
	ID int32
	// defines the name of the account e.g. equinix
	Name string
	Data map[project.Name]project.Data
)

// internal state
type Topology struct {
	Name             *Name             `json:"name"`
	Projects         Projects          `json:"projects"`
	CreatedAt        time.Time         `json:"created_at"`
	ProjectIDCounter *int32            `json:"project_id_counter"`
	Config           *staticres.Config `json:"config"`
}

type Projects map[project.ID]*project.Topology

var NotFoundID ID = -1
