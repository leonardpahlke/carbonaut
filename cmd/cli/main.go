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

package main

import (
	"os"

	"log/slog"

	"carbonaut.cloud/internal/connector"
	"carbonaut.cloud/internal/connector/providers"
	"carbonaut.cloud/internal/core"
	"carbonaut.cloud/internal/metrics"
	"carbonaut.cloud/pkg/schema"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

func main() {
	slog.Info("Create a new Carbonaut Config")
	cfg := core.Configuration{
		Kind: "carbonaut",
		Meta: core.Meta{
			Name: "carbonaut",
		},
		Spec: core.Spec{
			Provider: []connector.Config{{
				Kind: connector.ENERGY_PROVIDER,
				Meta: connector.Meta{
					Name:     "scaphandre-1",
					Provider: providers.ENERGY_PROVIDER_SCAPHANDRE,
				},
				Spec: providers.EnergyProviderConfig{
					Endpoints: []string{"http://localhost:8080/metrics"},
				},
			}, {
				Kind: connector.GEOLOCATION_PROVIDER,
				Meta: connector.Meta{
					Name:     "ipapi-1",
					Provider: providers.GEO_PROVIDER_IPAPI,
				},
				Spec: providers.GeoProviderConfig{
					IPAddresses: []string{"8.8.8.8"},
				},
			}, {
				Kind: connector.EMISSION_PROVIDER,
				Meta: connector.Meta{
					Name:     "energy-map-1",
					Provider: providers.EMISSION_PROVIDER_ENERGY_MAP,
				},
				Spec: providers.EmissionProviderConfig{
					Geolocation: &schema.Geolocation{
						Region: "DE",
					},
					Energy: []*schema.Energy{{
						Amount: 2669267,
						Unit:   schema.MICROWATT,
						Name:   "Host",
					}},
				},
			}},
			Metrics: metrics.Config{
				MetricsPort:   8082,
				CollectorName: "carbonaut",
			},
		},
	}
	if err := defaults.Set(&cfg); err != nil {
		slog.Error(err.Error(), "failed to set defaults")
	}
	y, err := yaml.Marshal(&cfg)
	if err != nil {
		slog.Error(err.Error(), "failed to marshal config")
	}
	// write yaml to file
	slog.Info("Write config to file")
	if err := os.WriteFile("config.yaml", y, 0o600); err != nil {
		slog.Error(err.Error(), "failed to write config to file")
	}
	slog.Info("Done, file written to config.yaml")
}
