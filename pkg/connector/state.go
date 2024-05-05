package connector

import (
	"time"

	"carbonaut.dev/pkg/schema"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
	"carbonaut.dev/pkg/schema/provider/resources"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

type state struct {
	Accounts map[resources.AccountIdentifier]Account
}

func newState() *state {
	return &state{
		Accounts: map[resources.AccountIdentifier]Account{},
	}
}

type Account struct {
	Status                    Status
	Meta                      Meta
	DynamicResourceCollectors map[resources.ResourceIdentifier]ResourceState
}

type ResourceState struct {
	StaticResourceData    staticres.Data
	StaticEnvironmentData staticenv.Data
	Meta                  Meta
}

// META
type Meta struct {
	Plugin    schema.PluginName
	CreatedAt time.Time
	DeletedAt time.Time
	Err       []error
}

type Err struct {
	E error
	T time.Time
}

// STATUS
type Status string

const (
	ToCreate Status = "to-create"
	Created  Status = "created"
)

// resourceReference is used internally to transport information about newly discovered resources within an account
type resourceReference struct {
	accountIdentifier  resources.AccountIdentifier
	resourceIdentifier resources.ResourceIdentifier
}
