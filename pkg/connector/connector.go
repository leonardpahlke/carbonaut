package connector

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"carbonaut.dev/pkg/connector/state"
	"carbonaut.dev/pkg/provider"
	"carbonaut.dev/pkg/provider/resource"
	"carbonaut.dev/pkg/provider/topology"
	"carbonaut.dev/pkg/util/compareutils"
)

type C struct {
	mutex           sync.Mutex
	connectorConfig *Config
	providerConfig  *provider.Config
	state           *state.S
}

type Config struct {
	TimeoutSeconds int `default:"60" json:"timeout_seconds" yaml:"timeout_seconds"`
}

func New(connectorConfig *Config, providerConfig *provider.Config) (*C, error) {
	connector := C{
		mutex:           sync.Mutex{},
		connectorConfig: connectorConfig,
		providerConfig:  &provider.Config{},
		state:           state.New(),
	}
	if err := connector.LoadConfig(providerConfig); err != nil {
		return nil, err
	}
	return &connector, nil
}

func (c *C) LoadConfig(newConfig *provider.Config) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var currentAccountSet, newAccountSet []*resource.AccountName

	buildAccountSet := func(resources provider.ResConfig) []*resource.AccountName {
		accountSet := make([]*resource.AccountName, 0, len(resources))
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

	slog.Debug("new carbonaut configuration parsed",
		"component", "connector.LoadConfig",
		"unaltered accounts", remainingAccounts,
		"deleted accounts", toBeDeletedAccounts,
		"new accounts", toBeCreatedAccounts,
	)

	// INFO: remainingAccounts are already configured and therefore no changes need to be made to the state

	// remove toBeDeletedAccounts from state
	for i := range toBeDeletedAccounts {
		slog.Debug("delete account from carbonaut state", "identifier", string(*toBeDeletedAccounts[i]))
		c.state.RemoveAccount(c.state.GetAccountID(toBeDeletedAccounts[i]))
	}

	// add toBeCreatedAccounts to "to-create" in state
	for i := range toBeCreatedAccounts {
		slog.Debug("added account to carbonaut state", "identifier", toBeCreatedAccounts[i])
		c.state.AddAccount(&topology.AccountT{
			Name:             toBeCreatedAccounts[i],
			Projects:         make(map[topology.ProjectID]*topology.ProjectT),
			CreatedAt:        time.Now(),
			ProjectIDCounter: new(int32),
			Config:           newConfig.Resources[*toBeCreatedAccounts[i]].StaticResConfig,
		})
	}

	slog.Info("configuration applied")
	c.providerConfig = newConfig

	return nil
}

// This function is run by the main control loop concurrently
func (c *C) Run(stopChan chan int, errChan chan error) {
	go func() {
		for {
			c.mutex.Lock()
			slog.Debug("start connector Run cycle")
			for aID := range c.state.T.Accounts {
				if err := c.updateStaticData(&aID); err != nil {
					errMsg := fmt.Errorf("unable to fetch resources, err: %v", err)
					slog.Error("unable to fetch resources", "error", err)
					errChan <- errMsg
				}
			}

			c.mutex.Unlock()
			slog.Debug("finished connector Run cycle")
			time.Sleep(time.Duration(c.connectorConfig.TimeoutSeconds) * time.Second)
		}
	}()
	<-stopChan
	slog.Info("received signal to stop the connector, shutting down")
}

func (c *C) GetStaticData() *topology.T {
	s := c.state.T
	return &s
}
