package connector

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"carbonaut.dev/pkg/plugin/staticenvplugin"
	"carbonaut.dev/pkg/plugin/staticresplugin"
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/environment/dynenv"
	"carbonaut.dev/pkg/schema/provider/environment/staticenv"
	"carbonaut.dev/pkg/schema/provider/resources/dynres"
	"carbonaut.dev/pkg/util/compareutils"
)

type C struct {
	mutex           sync.Mutex
	connectorConfig *Config
	providerConfig  *provider.Config
	state           *state
}

type Config struct {
	TimeoutSeconds int `json:"timeout_seconds"`
	Log            *slog.Logger
}

func New(connectorConfig *Config, providerConfig *provider.Config) (*C, error) {
	connector := C{
		mutex:           sync.Mutex{},
		connectorConfig: connectorConfig,
		providerConfig:  &provider.Config{},
		state:           newState(),
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

	c.connectorConfig.Log.Debug("new carbonaut configuration parsed",
		"unaltered accounts", remainingAccounts,
		"deleted accounts", toBeDeletedAccounts,
		"new accounts", toBeCreatedAccounts,
	)

	// INFO: remainingAccounts are already configured and therefore no changes need to be made to the state

	// remove toBeDeletedAccounts from state
	for i := range toBeDeletedAccounts {
		c.connectorConfig.Log.Debug("delete account from carbonaut state", "identifier", string(toBeDeletedAccounts[i]))
		delete(c.state.Accounts, toBeDeletedAccounts[i])
	}

	// add toBeCreatedAccounts to "to-create" in state
	for i := range toBeCreatedAccounts {
		c.connectorConfig.Log.Debug("added account to carbonaut state", "identifier", toBeCreatedAccounts[i])
		c.state.Accounts[toBeCreatedAccounts[i]] = Account{
			Meta: Meta{
				Plugin:    newConfig.Resources[toBeCreatedAccounts[i]].StaticResource.Plugin,
				CreatedAt: time.Now(),
			},
			DiscoveredResources: map[plugin.ResourceIdentifier]ResourceState{},
		}
	}

	c.connectorConfig.Log.Info("configuration applied")
	c.providerConfig = newConfig

	return nil
}

// This function is run by the main control loop concurrently
func (c *C) Run(stopChan chan int, errChan chan error) {
	go func() {
		for {
			c.mutex.Lock()
			c.connectorConfig.Log.Debug("start connector Run cycle")
			for accountIdentifier := range c.state.Accounts {
				toBeDeletedResources, toBeCreatedResources, err := c.fetchRemoteResourceState(accountIdentifier)
				if err != nil {
					errMsg := fmt.Errorf("unable to fetch resources, err: %v", err)
					c.connectorConfig.Log.Error("error", errMsg)
					errChan <- errMsg
				}

				if err := c.updateLocalResourceState(accountIdentifier, toBeDeletedResources, toBeCreatedResources); err != nil {
					errMsg := fmt.Errorf("unable to update resource data, err: %v", err)
					c.connectorConfig.Log.Error("error", errMsg)
					errChan <- errMsg
				}
			}

			c.mutex.Unlock()
			c.connectorConfig.Log.Debug("finished connector Run cycle")
			time.Sleep(time.Duration(c.connectorConfig.TimeoutSeconds) * time.Second)
		}
	}()
	<-stopChan
	c.connectorConfig.Log.Debug("received signal to stop the connector, shutting down")
}

func (c *C) fetchRemoteResourceState(accountIdentifier plugin.AccountIdentifier) ([]plugin.ResourceIdentifier, []plugin.ResourceIdentifier, error) {
	c.connectorConfig.Log.Debug("fetch resources", "account identifier", accountIdentifier)
	staticResPlugin := c.state.Accounts[accountIdentifier].Meta.Plugin
	p, err := staticresplugin.GetPlugin(staticResPlugin)
	if err != nil {
		c.connectorConfig.Log.Error("unable to find plugin", "provider type", "staticresplugin", "plugin identifier", staticResPlugin, "error", err)
		return nil, nil, err
	}
	disoveredResources, err := p.ListResources(c.providerConfig.Resources[accountIdentifier].StaticResource)
	if err != nil {
		c.connectorConfig.Log.Error("unable to list ressources", "provider type", "staticresplugin", "plugin identifier", staticResPlugin, "error", err)
		return nil, nil, err
	}

	currentResources := make([]plugin.ResourceIdentifier, 0, len(c.state.Accounts[accountIdentifier].DiscoveredResources))

	for res := range c.state.Accounts[accountIdentifier].DiscoveredResources {
		currentResources = append(currentResources, res)
	}

	remainingResources, toBeDeletedResources, toBeCreatedResources := compareutils.CompareLists(disoveredResources, currentResources)

	c.connectorConfig.Log.Debug("resources discovered",
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
		c.connectorConfig.Log.Debug("delete resource", "account identifier", accountIdentifier, "resource identifier", string(toBeDeletedResources[i]))
		delete(updatedAccountDetails.DiscoveredResources, toBeDeletedResources[i])
	}

	for i := range toBeCreatedResources {
		c.connectorConfig.Log.Debug("add resource to state", "account identifier", accountIdentifier, "resource identifier", toBeCreatedResources[i])

		staticResCollector, err := staticresplugin.GetPlugin(c.providerConfig.Resources[accountIdentifier].StaticResource.Plugin)
		if err != nil {
			return err
		}

		staticEnvCollector, err := staticenvplugin.GetPlugin(c.providerConfig.Environment.StaticEnvironment.Plugin)
		if err != nil {
			return err
		}

		c.connectorConfig.Log.Debug("request static resource information", "account identifier", accountIdentifier, "resource identifier", toBeCreatedResources[i])
		staticResData, err := staticResCollector.GetResource(c.providerConfig.Resources[accountIdentifier].StaticResource, toBeCreatedResources[i])
		if err != nil {
			return err
		}

		c.connectorConfig.Log.Debug("request static environment information", "account identifier", accountIdentifier, "resource identifier", toBeCreatedResources[i])
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

	c.connectorConfig.Log.Debug("add new resources to account")
	c.state.Accounts[accountIdentifier] = updatedAccountDetails

	return nil
}

// This function is triggered by the user interface
func (c *C) Collect() (*provider.Data, error) {
	c.mutex.Lock()
	c.connectorConfig.Log.Debug("collect data")
	data := make(provider.Data)

	for accountIdentifier := range c.state.Accounts {
		dataAccount := []provider.AccountData{}
		for _, resourceDetails := range c.state.Accounts[accountIdentifier].DiscoveredResources {
			// collect dynamic data - environment
			// TODO dynresplugin.GetPlugin()

			// collect dynamic data - resource
			// TODO: dynenvplugin.GetPlugin()

			dataAccount = append(dataAccount, provider.AccountData{
				StaticResourceData:     resourceDetails.StaticResourceData,
				DynamicResourceData:    dynres.Data{},
				StaticEnvironmentData:  resourceDetails.StaticEnvironmentData,
				DynamicEnvironmentData: dynenv.Data{},
			})
		}
		data[accountIdentifier] = dataAccount
	}

	c.mutex.Unlock()
	return &data, nil
}
