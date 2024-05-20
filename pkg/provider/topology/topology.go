package topology

import (
	"time"

	"carbonaut.dev/pkg/provider/plugin"
	"carbonaut.dev/pkg/provider/resource"
	"carbonaut.dev/pkg/provider/types/staticres"
)

var (
	AccountNotFoundID  AccountID  = -1
	ProjectNotFoundID  ProjectID  = -1
	ResourceNotFoundID ResourceID = -1
)

type (
	AccountID  int32
	ProjectID  int32
	ResourceID int32
)

type T struct {
	Accounts          map[AccountID]*AccountT `json:"accounts" yaml:"accounts"`
	AccountsIDCounter *int32                  `json:"-"`
}

type AccountT struct {
	Name             *resource.AccountName `json:"name"`
	Projects         Projects              `json:"projects"`
	CreatedAt        time.Time             `json:"created_at"`
	ProjectIDCounter *int32                `json:"-"`
	Config           *staticres.Config     `json:"-"`
}

type Projects map[ProjectID]*ProjectT

type ProjectT struct {
	Name              *resource.ProjectName `json:"name"`
	Resources         Resources             `json:"resources"`
	CreatedAt         time.Time             `json:"created_at"`
	ResourceIDCounter *int32                `json:"-"`
}

type Resources map[ResourceID]*ResourceT

// internal state
type ResourceT struct {
	Name       *resource.ResourceName  `json:"name"`
	StaticData *resource.StaticResData `json:"static_data"`
	CreatedAt  time.Time               `json:"created_at"`
	Plugin     *plugin.Kind            `json:"plugin"`
}
