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
	TimeoutSeconds int `default:"60" json:"timeout_seconds" yaml:"timeout_seconds"`
}

func New(connectorConfig *Config, logger *slog.Logger, providerConfig *provider.Config) (*C, error) {
	connector := C{
		mutex:           sync.Mutex{},
		connectorConfig: connectorConfig,
		providerConfig:  &provider.Config{},
		state:           newState(),
		log:             logger,
	}
	if err := connector.LoadConfig(providerConfig); err != nil {
		return nil, err
	}
	return &connector, nil
}

func (c *C) LoadConfig(newConfig *provider.Config) error {
	var newAccountLen int
	var currentAccountLen int

	if c.providerConfig != nil && c.providerConfig.Resources != nil {
		currentAccountLen = len(*c.providerConfig.Resources)
	}
	if newConfig != nil && newConfig.Resources != nil {
		newAccountLen = len(*newConfig.Resources)
	}

	newAccountSet := make([]plugin.AccountIdentifier, 0, newAccountLen)
	currentAccountSet := make([]plugin.AccountIdentifier, 0, currentAccountLen)

	if c.providerConfig != nil && c.providerConfig.Resources != nil {
		for k := range *c.providerConfig.Resources {
			currentAccountSet = append(currentAccountSet, k)
		}
	}

	if newConfig != nil && newConfig.Resources != nil {
		for k := range *newConfig.Resources {
			newAccountSet = append(newAccountSet, k)
		}
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
		if resConfig, ok := (*newConfig.Resources)[toBeCreatedAccounts[i]]; ok {
			c.state.Accounts[toBeCreatedAccounts[i]] = Account{
				Meta: Meta{
					Plugin:    resConfig.StaticResource.Plugin,
					CreatedAt: time.Now(),
				},
				DiscoveredResources: map[plugin.ResourceIdentifier]ResourceState{},
			}
		} else {
			c.log.Error("resource config not found for account", "account", toBeCreatedAccounts[i])
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

func (c *C) fetchRemoteResourceState(accountIdentifier plugin.AccountIdentifier) (toBeDeleted, toBeCreated []plugin.ResourceIdentifier, err error) {
	c.log.Info("fetch resources", "component", "connector.fetchRemoteResourceState", "account identifier", accountIdentifier)
	p, err := staticresplugin.GetPlugin(c.state.Accounts[accountIdentifier].Meta.Plugin)
	if err != nil {
		c.log.Error("unable to find plugin", "component", "connector.fetchRemoteResourceState", "provider type", "staticresplugin", "error", err)
		return nil, nil, err
	}

	var discoveredResources *[]plugin.ResourceIdentifier
	if staticAccountResources, ok := (*c.providerConfig.Resources)[accountIdentifier]; ok {
		discoveredResources, err = p.ListResources(staticAccountResources.StaticResource)
		if err != nil {
			c.log.Error("unable to list resources", "component", "connector.fetchRemoteResourceState", "provider type", "staticresplugin", "error", err)
			return nil, nil, err
		}
	} else {
		c.log.Error("resource config not found in account", "account", accountIdentifier)
	}

	currentResources := make([]plugin.ResourceIdentifier, 0, len(c.state.Accounts[accountIdentifier].DiscoveredResources))
	for res := range c.state.Accounts[accountIdentifier].DiscoveredResources {
		currentResources = append(currentResources, res)
	}

	remainingResources, toBeDeleted, toBeCreated := compareutils.CompareLists(*discoveredResources, currentResources)

	c.log.Debug("resources discovered",
		"component", "connector.fetchRemoteResourceState",
		"provider type", "staticresplugin",
		"account identifier", accountIdentifier,
		"unaltered resources", remainingResources,
		"deleted resources", toBeDeleted,
		"new resources", toBeCreated,
	)

	return toBeDeleted, toBeCreated, nil
}

func (c *C) updateLocalResourceState(accountIdentifier plugin.AccountIdentifier, toBeDeletedResources, toBeCreatedResources []plugin.ResourceIdentifier) error {
	updatedAccountDetails := c.state.Accounts[accountIdentifier]

	for i := range toBeDeletedResources {
		delete(updatedAccountDetails.DiscoveredResources, toBeDeletedResources[i])
	}

	if staticAccountResources, ok := (*c.providerConfig.Resources)[accountIdentifier]; ok {
		for i := range toBeCreatedResources {
			staticResCollector, err := staticresplugin.GetPlugin(staticAccountResources.StaticResource.Plugin)
			if err != nil {
				return err
			}

			staticEnvCollector, err := staticenvplugin.GetPlugin(c.providerConfig.Environment.StaticEnvironment.Plugin)
			if err != nil {
				return err
			}

			staticResData, err := staticResCollector.GetResource(staticAccountResources.StaticResource, &toBeCreatedResources[i])
			if err != nil {
				return err
			}

			staticEnvData, err := staticEnvCollector.Get(c.providerConfig.Environment.StaticEnvironment, &staticenv.InfraData{
				IP: staticResData.IP,
			})
			if err != nil {
				return err
			}

			updatedAccountDetails.DiscoveredResources[toBeCreatedResources[i]] = ResourceState{
				StaticResourceData:    staticResData,
				StaticEnvironmentData: staticEnvData,
				Meta: Meta{
					Plugin:    staticAccountResources.DynamicResource.Plugin,
					CreatedAt: time.Now(),
				},
			}
		}
	} else {
		c.log.Error("resource config not found in account", "account", accountIdentifier)
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
		if staticAccountResources, ok := (*c.providerConfig.Resources)[accountIdentifier]; ok {
			c.log.Debug("collect data", "component", "connector.collect", "account", accountIdentifier)
			dataAccount := []provider.AccountData{}
			for k := range c.state.Accounts[accountIdentifier].DiscoveredResources {
				c.log.Debug("collect dynamic data - resource", "component", "connector.collect", "plugin", *staticAccountResources.DynamicResource.Plugin)
				dynResCollector, err := dynresplugin.GetPlugin(staticAccountResources.DynamicResource.Plugin)
				if err != nil {
					return nil, err
				}

				c.log.Debug("collect dynamic data - environment", "component", "connector.collect", "plugin", *c.providerConfig.Environment.DynamicEnvironment.Plugin)
				dynEnvCollector, err := dynenvplugin.GetPlugin(c.providerConfig.Environment.DynamicEnvironment.Plugin)
				if err != nil {
					return nil, err
				}

				dynResData, err := dynResCollector.Get(staticAccountResources.DynamicResource, c.state.Accounts[accountIdentifier].DiscoveredResources[k].StaticResourceData)
				if err != nil {
					return nil, err
				}
				dynEnvData, err := dynEnvCollector.Get(c.providerConfig.Environment.DynamicEnvironment, c.state.Accounts[accountIdentifier].DiscoveredResources[k].StaticEnvironmentData)
				if err != nil {
					return nil, err
				}

				dataAccount = append(dataAccount, provider.AccountData{
					StaticResourceData:     c.state.Accounts[accountIdentifier].DiscoveredResources[k].StaticResourceData,
					DynamicResourceData:    dynResData,
					StaticEnvironmentData:  c.state.Accounts[accountIdentifier].DiscoveredResources[k].StaticEnvironmentData,
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
	}

	return &data, nil
}
