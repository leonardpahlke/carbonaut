/*
Copyright 2023 CARBONAUT AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package connector

import (
	"fmt"
	"log/slog"

	"carbonaut.cloud/internal/connector/providers"
	"carbonaut.cloud/pkg/schema"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
)

type IProvider interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
	ListenOnChannel()
}

type Config struct {
	// APIVersion string `yaml:"apiVersion" json:"apiVersion" validate:"required" default:"v1_alpha1"`
	// Provider kind like EnergyProvider
	Kind ProviderKinds `yaml:"kind" json:"kind" validate:"required"`
	Meta Meta          `yaml:"meta" validate:"required" json:"meta"`
	// Provider specific spec which gets resolved and casted during runtime
	Spec any `yaml:"spec" validate:"required" json:"spec"`
}

type Meta struct {
	Name     string `yaml:"name" json:"name" validate:"required"`
	Provider string `yaml:"provider" json:"provider" validate:"required"`
}

type ProviderKinds string

const (
	ENERGY_PROVIDER      ProviderKinds = "energy-provider"
	EMISSION_PROVIDER    ProviderKinds = "emission-provider"
	ENERGY_MIX_PROVIDER  ProviderKinds = "energy-mix-provider"
	GEOLOCATION_PROVIDER ProviderKinds = "geolocation-provider"
	IT_RESOURCE_PROVIDER ProviderKinds = "it-resource-provider"
)

func CreateProviders(cfg []Config, collectorName string) ([]IProvider, error) {
	slog.Info("create providers...")
	var configuredProviders []IProvider
	inChannelITResources := make(chan *schema.ITResource)
	inChannelEnergy := make(chan *schema.Energy)
	inChannelGeolocation := make(chan *schema.Geolocation)

	var multiErr error
	for i := range cfg {
		slog.Info("create new provider", slog.String("p", string(cfg[i].Kind)))
		switch cfg[i].Kind {
		case ENERGY_PROVIDER:
			energyP, err := providers.NewEnergyProvider(cfg[i].Spec, string(cfg[i].Meta.Provider), collectorName, inChannelITResources)
			if err != nil {
				multiErr = multierror.Append(multiErr, fmt.Errorf("error while initializing energy provider: %s", err.Error()))
				continue
			}
			configuredProviders = append(configuredProviders, energyP)
		case GEOLOCATION_PROVIDER:
			geoP, err := providers.NewGeoProvider(cfg[i].Spec, string(cfg[i].Meta.Provider), collectorName, inChannelITResources)
			if err != nil {
				multiErr = multierror.Append(multiErr, fmt.Errorf("error while initializing geolocation provider: %s", err.Error()))
				continue
			}
			configuredProviders = append(configuredProviders, geoP)
		case EMISSION_PROVIDER:
			emissionP, err := providers.NewEmissionProvider(cfg[i].Spec, string(cfg[i].Meta.Provider), collectorName, inChannelEnergy, inChannelGeolocation)
			if err != nil {
				multiErr = multierror.Append(multiErr, fmt.Errorf("error while initializing emission provider: %s", err.Error()))
				continue
			}
			configuredProviders = append(configuredProviders, emissionP)
			// TODO: add other providers
		default:
			multiErr = multierror.Append(multiErr, fmt.Errorf("unknown provider kind %s", cfg[i].Kind))
		}
	}
	if multiErr != nil {
		return nil, multiErr
	}

	slog.Info("✅ provider configured", slog.Int("n", len(configuredProviders)))
	return configuredProviders, nil
}
