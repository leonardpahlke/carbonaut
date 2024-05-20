package electricitymaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"carbonaut.dev/pkg/provider/account/project/resource"
	"carbonaut.dev/pkg/provider/environment"
	"carbonaut.dev/pkg/provider/plugin"
	"carbonaut.dev/pkg/provider/types/dynenv"
	"carbonaut.dev/pkg/util/cache"
	"carbonaut.dev/pkg/util/httpwrapper"
	"go.uber.org/multierr"
)

var (
	PluginName                    plugin.Kind = "electricitymaps"
	cacheZoneKey                  string      = "cache-zone-key"
	zoneLookupCacheTimeoutMinutes             = 60
)

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
		setupErrors = multierr.Append(setupErrors, errors.New("electricitymaps access key environment variable is not set or empty"))
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

func (p p) GetDynamicEnvironmentData(data *resource.Location) (*environment.DynamicEnvData, error) {
	resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
		Method:  http.MethodGet,
		BaseURL: "https://api.electricitymap.org/v3/power-breakdown/latest?zone=" + data.Country,
		Headers: map[string]string{
			"auth-token": *p.accessKey,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received none OK response code: %d, response %s", resp.StatusCode, resp.Body)
	}

	var zoneEnergyBreakdown ZoneEnergyBreakdown
	err = json.Unmarshal(resp.Body, &zoneEnergyBreakdown)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return zoneEnergyBreakdown.Convert(), nil
}

func (p p) QueryZones() (*Zones, error) {
	if cachedZones, found := p.cache.Get(cacheZoneKey); found {
		zones, ok := cachedZones.(*Zones)
		if !ok {
			return nil, errors.New("cached value is not of type *Zones")
		}
		slog.Debug("loaded electricitymap zones from cache")
		return zones, nil
	}

	resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
		Method:  http.MethodGet,
		BaseURL: "https://api.electricitymap.org/v3/zones",
	})
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	var zones Zones
	err = json.Unmarshal(resp.Body, &zones)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling project list request: %v", err)
	}
	if err := p.cache.Add(cacheZoneKey, &zones, time.Duration(zoneLookupCacheTimeoutMinutes)*time.Minute); err != nil {
		slog.Error("could not store zone information in cache", "error", err)
	}
	return &zones, nil
}
