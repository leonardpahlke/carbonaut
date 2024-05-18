package connector

import (
	"fmt"
	"log/slog"

	"carbonaut.dev/pkg/plugin/dynenvplugins"
	"carbonaut.dev/pkg/plugin/dynresplugins"
	"carbonaut.dev/pkg/provider"
	"carbonaut.dev/pkg/provider/data/account"
	"carbonaut.dev/pkg/provider/data/account/project"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/types/dynres"
)

// This function is triggered by the user interface
func (c *C) Collect() (*provider.Data, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	slog.Info("start collecting data")
	data := make(provider.Data)

	for aID := range c.state.T.Accounts {
		accountData := make(account.Data)
		if staticAccountResources, ok := (c.providerConfig.Resources)[*c.state.T.Accounts[aID].Name]; ok {
			slog.Debug("collect data", "account", aID)

			for p := range c.state.T.Accounts[aID].Projects {
				slog.Debug("collect data", "account", aID, "project", p)
				projectData := make(project.Data)
				for r := range c.state.T.Accounts[aID].Projects[p].Resources {
					dynData, err := c.collectDynResData(c.state.T.Accounts[aID].Projects[p].Resources[r], staticAccountResources.DynamicResConfig)
					if err != nil {
						return nil, err
					}
					projectData[*c.state.T.Accounts[aID].Projects[p].Resources[r].Name] = &resource.Data{
						DynamicData: dynData,
						StaticData:  c.state.T.Accounts[aID].Projects[p].Resources[r].StaticData,
					}
				}
				accountData[*c.state.T.Accounts[aID].Projects[p].Name] = projectData
			}
		}
		data[*c.state.T.Accounts[aID].Name] = accountData
	}

	slog.Info("data collected")

	return &data, nil
}

func (c *C) collectDynResData(r *resource.Topology, dynResConfig *dynres.Config) (*resource.DynamicData, error) {
	pRes, err := dynresplugins.GetPlugin(dynResConfig)
	if err != nil {
		return nil, fmt.Errorf("error loading dynres plugin: %s, err: %v", *dynResConfig.Plugin, err)
	}

	pEnv, err := dynenvplugins.GetPlugin(c.providerConfig.Environment.DynamicEnvConfig)
	if err != nil {
		return nil, fmt.Errorf("error loading dynenv plugin: %s, err: %v", *c.providerConfig.Environment.DynamicEnvConfig.Plugin, err)
	}

	slog.Debug("collect dynamic data - resource", "plugin", *r.Plugin)

	dynResData, err := pRes.GetDynamicResourceData(dynResConfig, r.StaticData)
	if err != nil {
		return nil, err
	}

	slog.Debug("collect dynamic data - environment", "plugin", *c.providerConfig.Environment.DynamicEnvConfig.Plugin)
	dynEnvData, err := pEnv.GetDynamicEnvironmentData(r.StaticData.Location)
	if err != nil {
		return nil, err
	}
	return &resource.DynamicData{
		ResData: dynResData,
		EnvData: dynEnvData,
	}, nil
}
