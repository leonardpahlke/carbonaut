package equinixplugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/types/staticres"
	"carbonaut.dev/pkg/util/cache"
	"carbonaut.dev/pkg/util/httpwrapper"
)

var PluginName plugin.Kind = "equinixplugin"

type p struct {
	cache *cache.Cache
}

func New() p {
	// Create a cache with an expiration time of 60 seconds, and which
	// purges expired items every 5 minutes
	c := cache.New(60*time.Second, 5*time.Minute)
	return p{
		cache: c,
	}
}

func (p p) GetName() *plugin.Kind {
	return &PluginName
}

func (p p) DiscoverProjectIdentifiers(cfg *staticres.Config) ([]*project.Name, error) {
	if cfg.Plugin == nil {
		return nil, errors.New("plugin is not set")
	}
	if cfg.AccessKey == nil || *cfg.AccessKey == "" {
		return nil, errors.New("access key is not set or empty")
	}

	resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
		Method:  http.MethodGet,
		BaseURL: "https://api.equinix.com/metal/v1/projects",
		Headers: map[string]string{
			"X-Auth-Token": *cfg.AccessKey,
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
		pName := project.Name(pResp.Projects[i].Name)
		prj = append(prj, &pName)
	}
	return prj, nil
}

func (p p) DiscoverStaticResourceIdentifiers(cfg *staticres.Config, pName *project.Name) ([]*resource.Name, error) {
	resData, err := p.getResourceData(cfg, pName)
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

func (p p) GetStaticResourceData(cfg *staticres.Config, pName *project.Name, rName *resource.Name) (*resource.StaticResData, error) {
	resData, err := p.getResourceData(cfg, pName)
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

func (p p) getResourceData(cfg *staticres.Config, pName *project.Name) (map[resource.Name]*resource.StaticResData, error) {
	cachedProjectResources, found := p.cache.Get(string(*pName))
	if found {
		resources, ok := cachedProjectResources.(map[resource.Name]*resource.StaticResData)
		if !ok {
			return nil, errors.New("cached value is not of type map[resource.Name]*resource.StaticResData")
		}
		return resources, nil
	}

	fetchedProjectResources, err := p.fetchResourceData(cfg, pName)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch resources for project %s", *pName)
	}

	if err := p.cache.Add(string(*pName), fetchedProjectResources, cache.DefaultExpiration); err != nil {
		return nil, errors.New("unable to add fetched projects to internal cache this can have a high performance impact")
	}
	return fetchedProjectResources, nil
}

func (p) fetchResourceData(cfg *staticres.Config, pName *project.Name) (map[resource.Name]*resource.StaticResData, error) {
	resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
		Method:  http.MethodGet,
		BaseURL: fmt.Sprintf("https://api.equinix.com/metal/v1/projects/%s/devices", *pName),
		Headers: map[string]string{
			"X-Auth-Token": *cfg.AccessKey,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	var devices ProjectDevicesResponse
	err = json.Unmarshal(resp.Body, &devices)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	resourceData := make(map[resource.Name]*resource.StaticResData)
	for i := range devices.Devices {
		d := devices.Devices[i]
		resourceData[resource.Name(devices.Devices[i].ID)] = EquinixDataIntegration(&d)
	}
	return resourceData, nil
}
