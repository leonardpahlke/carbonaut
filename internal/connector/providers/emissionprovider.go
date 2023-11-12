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

package providers

import (
	"fmt"
	"log/slog"

	"carbonaut.cloud/pkg/providers/electricitymap/electricitymap"
	"carbonaut.cloud/pkg/schema"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"
)

// TODO: carbon intensity can be get cached

const (
	EMISSION_PROVIDER_ENERGY_MAP = "energy-map"
)

type EmissionProviderConfig struct {
	Geolocation *schema.Geolocation
	Energy      []*schema.Energy
}

type EmissionProvider interface {
	GetEmission(zoneIdentifier string) (*schema.CarbonBreakdown, error)
}

var implementedEmissionProviders = map[string]EmissionProvider{
	EMISSION_PROVIDER_ENERGY_MAP: electricitymap.Provider{},
}

type EmissionP struct {
	config               EmissionProviderConfig
	kind                 string
	collectorName        string
	metrics              []EmissionMetricDetails
	inChannelEnergy      chan *schema.Energy
	inChannelGeolocation chan *schema.Geolocation
	OutChannelEmission   chan *schema.Emission
}

func NewEmissionProvider(spec any, kind string, collectorName string, inChannelEnergy chan *schema.Energy, inChannelGeolocation chan *schema.Geolocation) (*EmissionP, error) {
	var cfg EmissionProviderConfig
	specBytes, err := yaml.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling emission provider config: %w", err)
	}
	if err := yaml.Unmarshal(specBytes, &cfg); err != nil {
		return nil, fmt.Errorf("error while unmarshalling emission provider config: %w", err)
	}
	if err := cfg.VerifyConfig(kind); err != nil {
		return nil, fmt.Errorf("error while verifying the provider config: %w", err)
	}
	slog.Info("emission provider config verified", slog.String("kind", kind))
	return &EmissionP{
		config:        cfg,
		kind:          kind,
		collectorName: collectorName,
		metrics: []EmissionMetricDetails{
			{
				Name:        "emission_total_co2_kilogram",
				Type:        "gauge",
				Description: "Total emissions of resources",
				Labels:      []string{"total"},
				Fn:          collectTotalEmissionFootprint,
				metricDesc:  &prometheus.Desc{},
			},
			// TODO: this should be per region
			{
				Name:        "carbon_intensity_gram_per_kilowatt_hour",
				Type:        "gauge",
				Description: "Carbon intensity in gCO2eq/kWh",
				Labels:      []string{"total"},
				Fn:          collectTotalEmissionFootprint,
				metricDesc:  &prometheus.Desc{},
			},
		},
		inChannelEnergy:      inChannelEnergy,
		inChannelGeolocation: inChannelGeolocation,
		OutChannelEmission:   make(chan *schema.Emission),
	}, nil
}

func (g EmissionP) ListenOnChannel() {
	for r := range g.inChannelEnergy {
		slog.Info("received energy", slog.String("name", r.Name))
	}
}

func (g EmissionProviderConfig) VerifyConfig(kind string) error {
	// check if provider is implemented
	if _, ok := implementedEmissionProviders[kind]; !ok {
		return fmt.Errorf("provider %s is not implemented", kind)
	}
	return nil
}

func (p EmissionP) Describe(ch chan<- *prometheus.Desc) {
	for i := range p.metrics {
		p.metrics[i].metricDesc = prometheus.NewDesc(
			prometheus.BuildFQName(p.collectorName, "", p.metrics[i].Name),
			p.metrics[i].Description,
			p.metrics[i].Labels,
			nil,
		)
		slog.Info("metric description registered for", slog.String("metric name", p.metrics[i].Name))
	}
}

func (p EmissionP) Collect(ch chan<- prometheus.Metric) {
	slog.Info("collecting emission metrics")
	emissions, err := p.requestEmission(p.config)
	if err != nil {
		slog.Error("error while getting emission", slog.String("error", err.Error()))
		return
	}
	// TODO: add this feature later
	// slog.Info("sending emissions to other providers")
	// for j := range emissions {
	// 	p.OutChannelEmission <- emissions[j]
	// }
	slog.Info("register metrics")
	for j := range p.metrics {
		ch <- prometheus.MustNewConstMetric(p.metrics[0].metricDesc, prometheus.GaugeValue, float64(p.metrics[j].Fn(emissions)), p.metrics[j].Labels...)
	}
}

func (p EmissionP) requestEmission(cfg EmissionProviderConfig) ([]*schema.Emission, error) {
	slog.Info("get emission", "provider", p.kind)
	provider := implementedEmissionProviders[p.kind]
	carbonBreakdown, err := provider.GetEmission(cfg.Geolocation.Region)
	if err != nil {
		return nil, fmt.Errorf("error while getting emission: %w", err)
	}
	var emissions []*schema.Emission
	for i := range cfg.Energy {
		slog.Info("energy in kilowatts", slog.Float64("energy", cfg.Energy[i].ConvertToKilowatt()))
		emissions = append(emissions, &schema.Emission{
			CarbonIntensity: cfg.Energy[i].ConvertToKilowatt() * float64(carbonBreakdown.CarbonIntensity),
		})
	}
	slog.Info("collected emissions")
	return emissions, nil
}

type EmissionMetricDetails struct {
	Name        string
	Type        string
	Description string
	Labels      []string
	Fn          func(e []*schema.Emission) float64
	metricDesc  *prometheus.Desc
}

func collectTotalEmissionFootprint(e []*schema.Emission) float64 {
	var total float64
	for i := range e {
		total += e[i].CarbonIntensity
	}
	return total
}
