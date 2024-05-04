package connector

import (
	"fmt"
	"sync"
	"time"

	"carbonaut.dev/pkg/connector/provider"
	"carbonaut.dev/pkg/connector/provider/resources"
	"carbonaut.dev/pkg/util/compareutils"
)

type C struct {
	m               sync.Mutex
	connectorConfig *Config
	providerConfig  *provider.Config
	state           *state
	updatedConfig   bool
}

type Config struct {
	TimeoutSeconds int `json:"timeout_seconds"`
}

func New(connectorConfig *Config, providerConfig *provider.Config) *C {
	connector := C{
		m:               sync.Mutex{},
		connectorConfig: connectorConfig,
		providerConfig:  &provider.Config{},
		state:           newState(),
		updatedConfig:   false,
	}
	connector.LoadConfig(providerConfig)
	return &connector
}

func (c *C) LoadConfig(newConfig *provider.Config) error {
	// set deleted account to "to-delete" in state
	newAccountSet := make([]resources.AccountIdentifier, 0, len(c.providerConfig.Resources))
	currentAccountSet := make([]resources.AccountIdentifier, 0, len(newConfig.Resources))

	for k := range newConfig.Resources {
		newAccountSet = append(newAccountSet, k)
	}

	for k := range c.providerConfig.Resources {
		currentAccountSet = append(currentAccountSet, k)
	}

	remainingAccounts, toBeDeletedAccounts, toBeCreatedAccounts := compareutils.CompareLists(newAccountSet, currentAccountSet)

	fmt.Printf("accounts that are not changing: %v\n", remainingAccounts)
	fmt.Printf("accounts that are getting deleted: %v\n", toBeDeletedAccounts)
	fmt.Printf("accounts that are getting created: %v\n", toBeCreatedAccounts)

	// INFO: remainingAccounts are already configured and therefore no changes need to be made to the state

	// remove toBeDeletedAccounts from state
	for i := range toBeDeletedAccounts {
		fmt.Printf("remove account %s from state\n", toBeDeletedAccounts[i])
		delete(c.state.Accounts, toBeDeletedAccounts[i])
	}

	// add toBeCreatedAccounts to "to-create" in state
	for i := range toBeCreatedAccounts {
		fmt.Printf("add account %s to state\n", toBeCreatedAccounts[i])
		c.state.Accounts[toBeCreatedAccounts[i]] = Account{
			Status: ToCreate,
			Meta: Meta{
				Plugin:    newConfig.Resources[toBeCreatedAccounts[i]].StaticResource.Plugin,
				CreatedAt: time.Now(),
			},
			// this information will be added in the Run control loop
			DynamicResourceCollectors: map[resources.ResourceIdentifier]ResourceState{},
		}
		delete(c.state.Accounts, toBeDeletedAccounts[i])
	}

	fmt.Println("configuration applied")
	c.providerConfig = newConfig

	return nil
}

// This function is run by the main control loop concurrently
func (c *C) Run() {
	for {
		c.m.Lock()
		// c.cleanupRemovedAccounts()
		c.bootstrapNewAccounts()
		c.discoverResources()

		c.m.Unlock()
		time.Sleep(time.Duration(c.connectorConfig.TimeoutSeconds) * time.Second)
	}
}

// func (c *C) cleanupRemovedAccounts() {
// 	// DELETE clean up
// 	// 1. get all acounts with "to-delete" state
// 	// 2. remove associated resource observers
// 	// 3. set status to "deleted"
// }

func (c *C) bootstrapNewAccounts() {
	// BOOTSTRAP ACCOUNT
	// 1. get all accounts with "to-create" state
	// 2. fetch static data for each created account (infra & environment)
	// 3. register dynamic observers
	// 4. set status to "created"
}

func (c *C) discoverResources() {
	// BOOTSTRAP RESOURCES
	// 1. get all static resources that are "to-create" state
	// 2. register dynamic observers
	// 3. set status to "created" for the resource
}

// This function is triggered by the user interface
func (c *C) Collect() {}
