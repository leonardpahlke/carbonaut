package mockenergymix

import (
	"errors"
	"os"
	"time"

	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/data/environment"
	"carbonaut.dev/pkg/provider/types/dynenv"
	"carbonaut.dev/pkg/util/cache"
	"go.uber.org/multierr"
)

var PluginName plugin.Kind = "mockenergymix"

type p struct {
	cfg       *dynenv.Config
	cache     *cache.Cache
	accessKey *string
}

func New(cfg *dynenv.Config) (p, error) {
	// Create a cache with an expiration time of 60 seconds, and which
	// purges expired items every 5 minutes
	c := cache.New(60*time.Second, 5*time.Minute)

	authKey := os.Getenv(*cfg.AccessKeyEnv)
	var setupErrors error
	if cfg.Plugin == nil {
		setupErrors = multierr.Append(setupErrors, errors.New("plugin is not set information"))
	}
	if authKey == "" {
		setupErrors = multierr.Append(setupErrors, errors.New("access key environment variable is not set or empty"))
	}
	if setupErrors != nil {
		return p{}, setupErrors
	}
	return p{
		cfg:       cfg,
		cache:     c,
		accessKey: &authKey,
	}, nil
}

func (p) GetName() *plugin.Kind {
	return &PluginName
}

func (p) GetDynamicEnvironmentData(cfg *dynenv.Config, data *resource.Location) (*environment.DynamicEnvData, error) {
	return &environment.DynamicEnvData{
		SolarPercentage:        20.5,
		WindPercentage:         30.0,
		HydroPercentage:        15.0,
		NuclearPercentage:      10.0,
		FossilFuelsPercentage:  20.0,
		BioEnergyPercentage:    4.5,
		OtherSourcesPercentage: 0.0,
	}, nil
}
