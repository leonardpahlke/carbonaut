package connector

import (
	"fmt"

	"carbonaut.dev/pkg/plugins/dynenvplugins"
	"carbonaut.dev/pkg/plugins/dynresplugins"
	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/data/account"
	"carbonaut.dev/pkg/schema/provider/data/account/project"
	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
	"carbonaut.dev/pkg/schema/provider/types/dynres"
)

// This function is triggered by the user interface
func (c *C) Collect() (*provider.Data, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.log.Info("start collecting data")
	data := make(provider.Data)

	for aID := range c.state.Accounts {
		accountData := make(account.Data)
		if staticAccountResources, ok := (c.providerConfig.Resources)[aID]; ok {
			c.log.Debug("collect data", "account", aID)

			for p := range c.state.Accounts[aID].Projects {
				c.log.Debug("collect data", "account", aID, "project", p)
				projectData := make(project.Data)
				for r := range c.state.Accounts[aID].Projects[p].Resources {
					dynData, err := c.collectDynResData(&aID, c.state.Accounts[aID].Projects[p].Resources[r], staticAccountResources.DynamicResConfig)
					if err != nil {
						return nil, err
					}
					projectData[r] = &resource.Data{
						DynamicData: dynData,
						StaticData:  c.state.Accounts[aID].Projects[p].Resources[r].StaticData,
					}
				}
				accountData[p] = projectData
			}

		}
		data[aID] = accountData
	}

	c.log.Info("data collected")

	return &data, nil
}

func (c *C) collectDynResData(aID *account.ID, r *resource.Topology, DynamicResConfig *dynres.Config) (*resource.DynamicData, error) {
	pRes, found := dynresplugins.GetPlugin(DynamicResConfig.Plugin)
	if !found {
		return nil, fmt.Errorf("could not find plugin: %s", *DynamicResConfig.Plugin)
	}

	pEnv, found := dynenvplugins.GetPlugin(c.providerConfig.Environment.DynamicEnvConfig.Plugin)
	if !found {
		return nil, fmt.Errorf("could not find plugin: %s", *DynamicResConfig.Plugin)
	}

	c.log.Debug("collect dynamic data - resource", "plugin", *r.Plugin)

	dynResData, err := pRes.GetDynamicResourceData(DynamicResConfig, r.StaticData)
	if err != nil {
		return nil, err
	}

	c.log.Debug("collect dynamic data - environment", "plugin", *c.providerConfig.Environment.DynamicEnvConfig.Plugin)
	dynEnvData, err := pEnv.GetDynamicEnvironmentData(c.providerConfig.Environment.DynamicEnvConfig, r.StaticData.Location)
	if err != nil {
		return nil, err
	}
	return &resource.DynamicData{
		ResData: dynResData,
		EnvData: dynEnvData,
	}, nil
}
