package connector

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"carbonaut.dev/pkg/connector/state"
	"carbonaut.dev/pkg/provider"
	"carbonaut.dev/pkg/provider/data/account"
	"carbonaut.dev/pkg/provider/data/account/project"
	"carbonaut.dev/pkg/util/compareutils"
)

type C struct {
	mutex           sync.Mutex
	connectorConfig *Config
	providerConfig  *provider.Config
	state           *state.S
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
		state:           state.New(),
		log:             logger,
	}
	if err := connector.LoadConfig(providerConfig); err != nil {
		return nil, err
	}
	return &connector, nil
}

func (c *C) LoadConfig(newConfig *provider.Config) error {
	var currentAccountSet, newAccountSet []*account.Name

	buildAccountSet := func(resources provider.ResConfig) []*account.Name {
		accountSet := make([]*account.Name, 0, len(resources))
		for r := range resources {
			accountName := r
			accountSet = append(accountSet, &accountName)
		}
		return accountSet
	}

	if c.providerConfig != nil && c.providerConfig.Resources != nil {
		currentAccountSet = buildAccountSet(c.providerConfig.Resources)
	}
	if newConfig != nil && newConfig.Resources != nil {
		newAccountSet = buildAccountSet(newConfig.Resources)
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
		c.log.Debug("delete account from carbonaut state", "identifier", string(*toBeDeletedAccounts[i]))
		c.state.RemoveAccount(c.state.GetAccountID(toBeDeletedAccounts[i]))
	}

	// add toBeCreatedAccounts to "to-create" in state
	for i := range toBeCreatedAccounts {
		c.log.Debug("added account to carbonaut state", "identifier", toBeCreatedAccounts[i])
		c.state.AddAccount(&account.Topology{
			Name:             toBeCreatedAccounts[i],
			Projects:         make(map[project.ID]*project.Topology),
			CreatedAt:        time.Now(),
			ProjectIDCounter: new(int32),
			Config:           newConfig.Resources[*toBeCreatedAccounts[i]].StaticResConfig,
		})
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
			c.log.Debug("start connector Run cycle")
			for aID := range c.state.T.Accounts {
				if err := c.updateStaticData(&aID); err != nil {
					errMsg := fmt.Errorf("unable to fetch resources, err: %v", err)
					c.log.Error("error", errMsg)
					errChan <- errMsg
				}
			}

			c.mutex.Unlock()
			c.log.Debug("finished connector Run cycle")
			time.Sleep(time.Duration(c.connectorConfig.TimeoutSeconds) * time.Second)
		}
	}()
	<-stopChan
	c.log.Debug("received signal to stop the connector, shutting down")
}
