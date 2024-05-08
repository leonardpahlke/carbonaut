package connector

import (
	"time"

	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
	"carbonaut.dev/pkg/schema/provider/resources/staticres"
)

type state struct {
	Accounts map[plugin.AccountIdentifier]Account
}

func newState() *state {
	return &state{
		Accounts: map[plugin.AccountIdentifier]Account{},
	}
}

type Account struct {
	Meta                Meta
	DiscoveredResources map[plugin.ResourceIdentifier]ResourceState
}

type ResourceState struct {
	StaticResourceData    *staticres.Data
	StaticEnvironmentData *staticenv.Data
	Meta                  Meta
}

// META
type Meta struct {
	Plugin    *plugin.Kind
	CreatedAt time.Time
	DeletedAt time.Time
}
