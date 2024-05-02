package main

import (
	"fmt"

	"carbonaut.dev/pkg/carbonautserver"
	"carbonaut.dev/pkg/config"
	"carbonaut.dev/pkg/infrastructureaccount"
)

func main() {
	config_path := "config.yaml"
	config, err := config.ReadConfig(config_path)
	if err != nil {
		panic(err) // TODO: fail
	}

	exitChan := make(chan int)
	// TODO: perphaps pass in some information about the deployment which can be accessed and updated by the resource plugin runners
	go carbonautserver.Listen(exitChan)

	infrastructureaccounts := []infrastructureaccount.ResourceAccount{}
	for i := range config.StaticResourceProviders {
		a, err := infrastructureaccount.New(config.StaticResourceProviders[i])
		if err != nil {
			panic(err) // TODO: error handling
		}
		infrastructureaccounts = append(infrastructureaccounts, a)
		go a.Observe()
		// go collector.Collect(config.StaticResourceProviders[i])
	}

	<-exitChan
	fmt.Println("Shutting down Carbonaut")
}
