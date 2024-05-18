package equinixplugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/types/staticres"
	"carbonaut.dev/pkg/util/cache"
	"carbonaut.dev/pkg/util/compareutils"
	"carbonaut.dev/pkg/util/httpwrapper"
	"go.uber.org/multierr"
)

var PluginName plugin.Kind = "equinix"

type p struct {
	cfg       *staticres.Config
	cache     *cache.Cache
	accessKey *string
}

func New(cfg *staticres.Config) (p, error) {
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

func (p p) GetName() *plugin.Kind {
	return &PluginName
}

func (p p) DiscoverProjectIdentifiers() ([]*project.Name, error) {
	resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
		Method:  http.MethodGet,
		BaseURL: "https://api.equinix.com/metal/v1/projects",
		Headers: map[string]string{
			"X-Auth-Token": *p.accessKey,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	var pResp ProjectsResponse
	err = json.Unmarshal(resp.Body, &pResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling project list request: %v", err)
	}
	prj := make([]*project.Name, 0)
	for i := range pResp.Projects {
		slog.Debug("discovered equinix project", "project name", pResp.Projects[i].Name)
		pName := project.Name(pResp.Projects[i].ID)
		prj = append(prj, &pName)
	}
	return prj, nil
}

func (p p) DiscoverStaticResourceIdentifiers(pName *project.Name) ([]*resource.Name, error) {
	resData, err := p.PGetResourceData(pName)
	if err != nil {
		return nil, fmt.Errorf("error getResourceData error: %v", err)
	}
	resourceNames := []*resource.Name{}
	for i := range resData {
		rName := i
		resourceNames = append(resourceNames, &rName)
	}
	return resourceNames, nil
}

func (p p) GetStaticResourceData(pName *project.Name, rName *resource.Name) (*resource.StaticResData, error) {
	resData, err := p.PGetResourceData(pName)
	if err != nil {
		return nil, fmt.Errorf("error getResourceData error: %v", err)
	}
	for i := range resData {
		if i == *rName {
			return resData[i], nil
		}
	}
	return nil, errors.New("resource not found")
}

func (p p) PGetResourceData(pName *project.Name) (map[resource.Name]*resource.StaticResData, error) {
	if cachedProjectResources, found := p.cache.Get(string(*pName)); found {
		resources, ok := cachedProjectResources.(map[resource.Name]*resource.StaticResData)
		if !ok {
			return nil, errors.New("cached value is not of type map[resource.Name]*resource.StaticResData")
		}
		slog.Debug("loaded equinix resources from cache", "resource names", compareutils.CollectKeys(resources))
		return resources, nil
	}

	fetchedProjectResources, err := p.PRequestResourceData(pName)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch resources for project %s", *pName)
	}

	if err := p.cache.Add(string(*pName), fetchedProjectResources, cache.DefaultExpiration); err != nil {
		return nil, errors.New("unable to add fetched projects to internal cache this can have a high performance impact")
	}
	return fetchedProjectResources, nil
}

func (p p) PRequestResourceData(pName *project.Name) (map[resource.Name]*resource.StaticResData, error) {
	resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
		Method:  http.MethodGet,
		BaseURL: fmt.Sprintf("https://api.equinix.com/metal/v1/projects/%s/devices", *pName),
		Headers: map[string]string{
			"X-Auth-Token": *p.accessKey,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	var devices EquinixDeviceResponse
	err = json.Unmarshal(resp.Body, &devices)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	resourceData := make(map[resource.Name]*resource.StaticResData)
	if len(devices.Devices) != 0 {
		for i := range devices.Devices {
			d := devices.Devices[i]
			resourceData[resource.Name(devices.Devices[i].ID)] = EquinixDataIntegration(&d)
		}
		slog.Debug("requested equinix resources", "resource names", compareutils.CollectKeys(resourceData))
	} else {
		slog.Info("no equinix resources found", "project name", *pName)
	}
	return resourceData, nil
}
