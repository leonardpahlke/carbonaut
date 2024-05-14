package equinixplugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/data/account/project"
	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
	"carbonaut.dev/pkg/schema/provider/types/staticres"
	"carbonaut.dev/pkg/util/httpwrapper"
)

var PluginName plugin.Kind = "equinixplugin"

type p struct{}

func New() p {
	return p{}
}

func (p) GetName() *plugin.Kind {
	return &PluginName
}

func (p) DiscoverProjectIdentifiers(cfg *staticres.Config) ([]*project.Name, error) {
	return nil, errors.New("not yet implemented")
}

func (p) DiscoverStaticResourceIdentifiers(cfg *staticres.Config, pName *project.Name) ([]*resource.Name, error) {
	return nil, errors.New("not yet implemented")
}

func (p) GetStaticResourceData(cfg *staticres.Config, pName *project.Name, rName *resource.Name) (*resource.StaticResData, error) {
	return nil, errors.New("not yet implemented")
}

func (p) GetStaticEnvironmentAndResourceData(cfg *staticres.Config) (*resource.StaticData, error) {
	if cfg.Plugin == nil {
		return nil, errors.New("[equinix-plugin] plugin is not set")
	}
	if cfg.AccessKey == nil || *cfg.AccessKey == "" {
		return nil, errors.New("[equinix-plugin] access key is not set or empty")
	}

	resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
		Method:  http.MethodGet,
		BaseURL: "https://api.equinix.com/metal/v1/projects",
		Headers: map[string]string{
			"X-Auth-Token": *cfg.AccessKey,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("[equinix-plugin] error creating request: %v", err)
	}

	var projects ProjectsResponse
	err = json.Unmarshal(resp.Body, &projects)
	if err != nil {
		return nil, fmt.Errorf("[equinix-plugin] error unmarshalling project list request: %v", err)
	}

	// data := make(staticenvres.Data)
	// for i := range projects.Projects {
	// 	projectDevices, err := fetchProjectResources(cfg, projects.Projects[i].ID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("[equinix-plugin] could not retrieve devices of project %s with ID %s. error: %v", projects.Projects[i].Name, projects.Projects[i].ID, err)
	// 	}
	// 	accountData := []staticenvres.StaticData{}
	// 	for j := range projectDevices.Devices {
	// 		staticresData, staticenvData := equinixDataIntegration(&projectDevices.Devices[j])
	// 		accountData = append(accountData, staticenvres.StaticData{
	// 			StaticResourceData:    staticresData,
	// 			StaticEnvironmentData: staticenvData,
	// 		})
	// 	}
	// 	data[account.ID(projects.Projects[i].ID)] = accountData
	// }
	// return &data, nil
	return nil, nil
}

// func fetchProjectResources(cfg *staticres.Config, projectID string) (*ProjectDevicesResponse, error) {
// 	resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
// 		Method:  http.MethodGet,
// 		BaseURL: fmt.Sprintf("https://api.equinix.com/metal/v1/projects/%s/devices", projectID),
// 		Headers: map[string]string{
// 			"X-Auth-Token": *cfg.AccessKey,
// 		},
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("[equinix-plugin] error creating request: %v", err)
// 	}

// 	var devices ProjectDevicesResponse
// 	err = json.Unmarshal(resp.Body, &devices)
// 	if err != nil {
// 		return nil, fmt.Errorf("[equinix-plugin] error reading response body: %v", err)
// 	}

// 	return &devices, nil
// }
