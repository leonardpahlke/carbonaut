package infrastructureaccount

import "carbonaut.dev/pkg/config"

type ResourceAccount struct {
	config config.StaticResourceConfig
}

func New(staticResourceConfig config.StaticResourceConfig) (ResourceAccount, error) {
	// TODO match based on plugin identifier
	return ResourceAccount{
		config: staticResourceConfig,
	}, nil
}

func (r ResourceAccount) Observe() {
	//
}
