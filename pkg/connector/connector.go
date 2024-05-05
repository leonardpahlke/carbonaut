package connector

import (
	"log/slog"
	"sync"
	"time"

	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/resources"
	"carbonaut.dev/pkg/util/compareutils"
)

type C struct {
	m              sync.Mutex
	cfg            *Config
	providerConfig *provider.Config
	state          *state
	configUpdated  bool
}

type Config struct {
	TimeoutSeconds int `json:"timeout_seconds"`
	Log            *slog.Logger
}

func New(connectorConfig *Config, providerConfig *provider.Config) (*C, error) {
	connector := C{
		m:              sync.Mutex{},
		cfg:            connectorConfig,
		providerConfig: &provider.Config{},
		state:          newState(),
		configUpdated:  false,
	}
	if err := connector.LoadConfig(providerConfig); err != nil {
		return nil, err
	}
	return &connector, nil
}

func (c *C) LoadConfig(newConfig *provider.Config) error {
	newAccountSet := make([]resources.AccountIdentifier, 0, len(c.providerConfig.Resources))
	currentAccountSet := make([]resources.AccountIdentifier, 0, len(newConfig.Resources))

	for k := range newConfig.Resources {
		newAccountSet = append(newAccountSet, k)
	}

	for k := range c.providerConfig.Resources {
		currentAccountSet = append(currentAccountSet, k)
	}

	remainingAccounts, toBeDeletedAccounts, toBeCreatedAccounts := compareutils.CompareLists(newAccountSet, currentAccountSet)

	c.cfg.Log.Debug("new carbonaut configuration parsed",
		"unaltered accounts", remainingAccounts,
		"deleted accounts", toBeDeletedAccounts,
		"new accounts", toBeCreatedAccounts,
	)

	// INFO: remainingAccounts are already configured and therefore no changes need to be made to the state

	// remove toBeDeletedAccounts from state
	for i := range toBeDeletedAccounts {
		c.cfg.Log.Debug("delete account from carbonaut state", "identifier", string(toBeDeletedAccounts[i]))
		delete(c.state.Accounts, toBeDeletedAccounts[i])
	}

	// add toBeCreatedAccounts to "to-create" in state
	for i := range toBeCreatedAccounts {
		c.cfg.Log.Debug("added account to carbonaut state", "identifier", toBeCreatedAccounts[i])
		c.state.Accounts[toBeCreatedAccounts[i]] = Account{
			Status: ToCreate,
			Meta: Meta{
				Plugin:    newConfig.Resources[toBeCreatedAccounts[i]].StaticResource.Plugin,
				CreatedAt: time.Now(),
			},
			DynamicResourceCollectors: map[resources.ResourceIdentifier]ResourceState{},
		}
	}

	c.configUpdated = true
	c.cfg.Log.Info("configuration applied")
	c.providerConfig = newConfig

	return nil
}

// This function is run by the main control loop concurrently
func (c *C) Run() {
	for {
		c.m.Lock()
		c.cfg.Log.Debug("start connector Run cycle")
		if c.configUpdated {
			c.cfg.Log.Debug("config updated, bootstrap new accounts")
			if err := c.bootstrapNewAccounts(); err != nil {
				c.cfg.Log.Error("unable to boostrap new accounts", "error", err)
			}
		}
		newDiscoveredResources, err := c.discoverResources()
		if err != nil {
			c.cfg.Log.Error("unable to discover static resource data", "error", err)
		}

		c.cfg.Log.Error("finished discovering new resources", "number of new resources discovered", len(newDiscoveredResources))

		if len(newDiscoveredResources) != 0 {
			if err := c.discoverEnvironment(newDiscoveredResources); err != nil {
				c.cfg.Log.Error("unable to discover static environment data", "error", err)
			}
		}

		c.m.Unlock()
		c.cfg.Log.Debug("finished connector Run cycle")
		time.Sleep(time.Duration(c.cfg.TimeoutSeconds) * time.Second)
	}
}

func (c *C) bootstrapNewAccounts() error {
	for k, account := range c.state.Accounts {
		if account.Status == ToCreate {
			c.cfg.Log.Debug("detected account which needs to get created", "account identifier", k)

			c.cfg.Log.Debug("fetch static resource data for account", "account identifier", k)

			// 3. register dynamic observers

			// 4. set status to "created"
			account.Status = Created
		}
	}

	return nil
}

func (c *C) discoverResources() ([]resourceReference, error) {
	// BOOTSTRAP RESOURCES
	// 1. get all static resources that are "to-create" state
	// 2. register dynamic observers
	// 3. set status to "created" for the resource
	return []resourceReference{}, nil
}

func (c *C) discoverEnvironment(resourceInfo []resourceReference) error {
	// BOOTSTRAP ENVIRONMENT DATA OF RESOURCES
	// 1. get all static resources that are "to-create" state
	// 2. register dynamic observers
	// 3. set status to "created" for the resource
	return nil
}

// This function is triggered by the user interface
func (c *C) Collect() {}
