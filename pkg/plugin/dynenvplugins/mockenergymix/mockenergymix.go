package mockenergymix

import (
	"errors"
	"os"
	"time"

	"carbonaut.dev/pkg/provider/environment"
	"carbonaut.dev/pkg/provider/plugin"
	"carbonaut.dev/pkg/provider/resource"
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
		setupErrors = multierr.Append(setupErrors, errors.New("mockenergymix access key environment variable is not set or empty"))
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

func (p) GetDynamicEnvironmentData(data *resource.Location) (*environment.DynamicEnvData, error) {
	return &environment.DynamicEnvData{
		SolarPercentage:        10,
		WindPercentage:         20,
		HydroPercentage:        10,
		NuclearPercentage:      20,
		GeothermalPercentage:   15,
		GasPercentage:          5,
		OilPercentage:          5,
		BiomassPercentage:      7,
		CoalPercentage:         8,
		OtherSourcesPercentage: 0,
		FossilFuelsPercentage:  38,
		RenewablePercentage:    62,
	}, nil
}
