package account

import (
	"time"

	"carbonaut.dev/pkg/schema/provider/data/account/project"
)

type (
	// defines the name of the account e.g. equinix
	ID   string
	Data map[project.ID]project.Data
)

type Topology struct {
	Projects  Projects  `json:"projects"         yaml:"projects"`
	CreatedAt time.Time `json:"created_at"         yaml:"created_at"`
}

type Projects map[project.ID]*project.Topology
