package connector

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"carbonaut.dev/pkg/plugin/dynenvplugin"
	"carbonaut.dev/pkg/plugin/dynresplugin"
	"carbonaut.dev/pkg/plugin/staticenvplugin"
	"carbonaut.dev/pkg/plugin/staticresplugin"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
	"carbonaut.dev/pkg/util/compareutils"
)

type C struct {
	mutex           sync.Mutex
	connectorConfig *Config
	providerConfig  *provider.Config
	state           *state
	log             *slog.Logger
}

type Config struct {
	TimeoutSeconds int `json:"timeout_seconds"`
}

func New(connectorConfig *Config, Log *slog.Logger, providerConfig *provider.Config) (*C, error) {
	connector := C{
		mutex:           sync.Mutex{},
		connectorConfig: connectorConfig,
		providerConfig:  &provider.Config{},
		state:           newState(),
		log:             Log,
	}
	if err := connector.LoadConfig(providerConfig); err != nil {
		return nil, err
	}
	return &connector, nil
}

func (c *C) LoadConfig(newConfig *provider.Config) error {
	newAccountSet := make([]plugin.AccountIdentifier, 0, len(c.providerConfig.Resources))
	currentAccountSet := make([]plugin.AccountIdentifier, 0, len(newConfig.Resources))

	for k := range newConfig.Resources {
		newAccountSet = append(newAccountSet, k)
	}

	for k := range c.providerConfig.Resources {
		currentAccountSet = append(currentAccountSet, k)
	}

	remainingAccounts, toBeDeletedAccounts, toBeCreatedAccounts := compareutils.CompareLists(newAccountSet, currentAccountSet)

	c.log.Debug("new carbonaut configuration parsed",
		"component", "connector.LoadConfig",
		"unaltered accounts", remainingAccounts,
		"deleted accounts", toBeDeletedAccounts,
		"new accounts", toBeCreatedAccounts,
	)

	// INFO: remainingAccounts are already configured and therefore no changes need to be made to the state

	// remove toBeDeletedAccounts from state
	for i := range toBeDeletedAccounts {
		c.log.Debug("delete account from carbonaut state", "component", "connector.LoadConfig", "identifier", string(toBeDeletedAccounts[i]))
		delete(c.state.Accounts, toBeDeletedAccounts[i])
	}

	// add toBeCreatedAccounts to "to-create" in state
	for i := range toBeCreatedAccounts {
		c.log.Debug("added account to carbonaut state", "component", "connector.LoadConfig", "identifier", toBeCreatedAccounts[i])
		c.state.Accounts[toBeCreatedAccounts[i]] = Account{
			Meta: Meta{
				Plugin:    newConfig.Resources[toBeCreatedAccounts[i]].StaticResource.Plugin,
				CreatedAt: time.Now(),
			},
			DiscoveredResources: map[plugin.ResourceIdentifier]ResourceState{},
		}
	}

	c.log.Info("configuration applied")
	c.providerConfig = newConfig

	return nil
}

// This function is run by the main control loop concurrently
func (c *C) Run(stopChan chan int, errChan chan error) {
	go func() {
		for {
			c.mutex.Lock()
			c.log.Debug("start connector Run cycle", "component", "connector.Run")
			for accountIdentifier := range c.state.Accounts {
				toBeDeletedResources, toBeCreatedResources, err := c.fetchRemoteResourceState(accountIdentifier)
				if err != nil {
					errMsg := fmt.Errorf("unable to fetch resources, err: %v", err)
					c.log.Error("error", errMsg)
					errChan <- errMsg
				}

				if err := c.updateLocalResourceState(accountIdentifier, toBeDeletedResources, toBeCreatedResources); err != nil {
					errMsg := fmt.Errorf("unable to update resource data, err: %v", err)
					c.log.Error("error", errMsg)
					errChan <- errMsg
				}
			}

			c.mutex.Unlock()
			c.log.Debug("finished connector Run cycle", "component", "connector.Run")
			time.Sleep(time.Duration(c.connectorConfig.TimeoutSeconds) * time.Second)
		}
	}()
	<-stopChan
	c.log.Debug("received signal to stop the connector, shutting down", "component", "connector.Run")
}

func (c *C) fetchRemoteResourceState(accountIdentifier plugin.AccountIdentifier) ([]plugin.ResourceIdentifier, []plugin.ResourceIdentifier, error) {
	c.log.Info("fetch resources", "component", "connector.fetchRemoteResourceState", "account identifier", accountIdentifier)
	staticResPlugin := c.state.Accounts[accountIdentifier].Meta.Plugin
	p, err := staticresplugin.GetPlugin(staticResPlugin)
	if err != nil {
		c.log.Error("unable to find plugin", "component", "connector.fetchRemoteResourceState", "provider type", "staticresplugin", "plugin identifier", staticResPlugin, "error", err)
		return nil, nil, err
	}
	disoveredResources, err := p.ListResources(c.providerConfig.Resources[accountIdentifier].StaticResource)
	if err != nil {
		c.log.Error("unable to list ressources", "component", "connector.fetchRemoteResourceState", "provider type", "staticresplugin", "plugin identifier", staticResPlugin, "error", err)
		return nil, nil, err
	}

	currentResources := make([]plugin.ResourceIdentifier, 0, len(c.state.Accounts[accountIdentifier].DiscoveredResources))

	for res := range c.state.Accounts[accountIdentifier].DiscoveredResources {
		currentResources = append(currentResources, res)
	}

	remainingResources, toBeDeletedResources, toBeCreatedResources := compareutils.CompareLists(disoveredResources, currentResources)

	c.log.Debug("resources discovered",
		"component", "connector.fetchRemoteResourceState",
		"provider type", "staticresplugin",
		"account identifier", accountIdentifier,
		"plugin identifier", staticResPlugin,
		"unaltered resources", remainingResources,
		"deleted resources", toBeDeletedResources,
		"new resources", toBeCreatedResources,
	)

	return toBeDeletedResources, toBeCreatedResources, nil
}

func (c *C) updateLocalResourceState(accountIdentifier plugin.AccountIdentifier, toBeDeletedResources []plugin.ResourceIdentifier, toBeCreatedResources []plugin.ResourceIdentifier) error {
	updatedAccountDetails := c.state.Accounts[accountIdentifier]
	for i := range toBeDeletedResources {
		// c.log.Debug("delete resource", "component", "connector.updateLocalResourceState", "account identifier", accountIdentifier, "resource identifier", string(toBeDeletedResources[i]))
		delete(updatedAccountDetails.DiscoveredResources, toBeDeletedResources[i])
	}

	for i := range toBeCreatedResources {
		// c.log.Debug("add resource to state", "component", "connector.updateLocalResourceState", "account identifier", accountIdentifier, "resource identifier", toBeCreatedResources[i])

		staticResCollector, err := staticresplugin.GetPlugin(c.providerConfig.Resources[accountIdentifier].StaticResource.Plugin)
		if err != nil {
			return err
		}

		staticEnvCollector, err := staticenvplugin.GetPlugin(c.providerConfig.Environment.StaticEnvironment.Plugin)
		if err != nil {
			return err
		}

		// c.log.Debug("request static resource information", "component", "connector.updateLocalResourceState", "account identifier", accountIdentifier, "resource identifier", toBeCreatedResources[i])
		staticResData, err := staticResCollector.GetResource(c.providerConfig.Resources[accountIdentifier].StaticResource, toBeCreatedResources[i])
		if err != nil {
			return err
		}

		// c.log.Debug("request static environment information", "component", "connector.updateLocalResourceState", "account identifier", accountIdentifier, "resource identifier", toBeCreatedResources[i])
		staticEnvData, err := staticEnvCollector.Get(c.providerConfig.Environment.StaticEnvironment, staticenv.InfraData{
			IP: staticResData.IP,
		})
		if err != nil {
			return err
		}

		updatedAccountDetails.DiscoveredResources[toBeCreatedResources[i]] = ResourceState{
			StaticResourceData:    staticResData,
			StaticEnvironmentData: staticEnvData,
			Meta: Meta{
				Plugin:    c.providerConfig.Resources[accountIdentifier].DynamicResource.Plugin,
				CreatedAt: time.Now(),
			},
		}
	}

	c.log.Info("add new resources to account", "component", "connector.updateLocalResourceState")
	c.state.Accounts[accountIdentifier] = updatedAccountDetails

	return nil
}

// This function is triggered by the user interface
func (c *C) Collect() (*provider.Data, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.log.Info("collect data", "component", "connector.collect")
	data := make(provider.Data)

	for accountIdentifier := range c.state.Accounts {
		c.log.Debug("collect data", "component", "connector.collect", "account", accountIdentifier)
		dataAccount := []provider.AccountData{}
		for _, resourceDetails := range c.state.Accounts[accountIdentifier].DiscoveredResources {
			c.log.Debug("collect dynamic data - resource", "component", "connector.collect", "plugin", c.providerConfig.Resources[accountIdentifier].DynamicResource.Plugin)
			dynResCollector, err := dynresplugin.GetPlugin(c.providerConfig.Resources[accountIdentifier].DynamicResource.Plugin)
			if err != nil {
				return nil, err
			}

			c.log.Debug("collect dynamic data - environment", "component", "connector.collect", "plugin", c.providerConfig.Environment.DynamicEnvironment.Plugin)
			dynEnvCollector, err := dynenvplugin.GetPlugin(c.providerConfig.Environment.DynamicEnvironment.Plugin)
			if err != nil {
				return nil, err
			}

			dynResData, err := dynResCollector.Get(c.providerConfig.Resources[accountIdentifier].DynamicResource, resourceDetails.StaticResourceData)
			if err != nil {
				return nil, err
			}
			dynEnvData, err := dynEnvCollector.Get(c.providerConfig.Environment.DynamicEnvironment, resourceDetails.StaticEnvironmentData)
			if err != nil {
				return nil, err
			}

			dataAccount = append(dataAccount, provider.AccountData{
				StaticResourceData:     resourceDetails.StaticResourceData,
				DynamicResourceData:    dynResData,
				StaticEnvironmentData:  resourceDetails.StaticEnvironmentData,
				DynamicEnvironmentData: dynEnvData,
			})
		}

		if len(dataAccount) == 0 {
			c.log.Debug("No data collected for account", "component", "connector.collect", "account", accountIdentifier)
		} else {
			data[accountIdentifier] = dataAccount
			c.log.Debug("Data collected for account", "component", "connector.collect", "account", accountIdentifier, "dataCount", len(dataAccount))
		}
	}

	return &data, nil
}
