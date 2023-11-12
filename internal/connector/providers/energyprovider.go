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
	"regexp"

	"log/slog"

	"carbonaut.cloud/pkg/providers/scaphandre"
	"carbonaut.cloud/pkg/schema"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"
)

const (
	ENERGY_PROVIDER_SCAPHANDRE = "scaphandre"
)

type EnergyProviderConfig struct {
	// Endpoint of the provider; needs to be a valid URL or IP address
	Endpoints []string `yaml:"endpoints" json:"endpoints"`
}

type EnergyProvider interface {
	GetEnergy(endpoint string) ([]*schema.Energy, error)
}

type EnergyP struct {
	config               EnergyProviderConfig
	kind                 string
	collectorName        string
	metrics              []EnergyMetricDetails
	inChannelITResources chan *schema.ITResource
	OutChannelEnergy     chan *schema.Energy
}

func NewEnergyProvider(spec any, kind string, collectorName string, inChannelITResources chan *schema.ITResource) (*EnergyP, error) {
	var cfg EnergyProviderConfig
	specBytes, err := yaml.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling energy provider config: %w", err)
	}
	if err := yaml.Unmarshal(specBytes, &cfg); err != nil {
		return nil, fmt.Errorf("error while unmarshalling energy provider config: %w", err)
	}
	if err := cfg.VerifyConfig(kind); err != nil {
		return nil, fmt.Errorf("error while verifying the provider config: %w", err)
	}
	slog.Info("energy provider config verified", slog.String("kind", kind))
	return &EnergyP{
		config: cfg,
		kind:   kind,
		metrics: []EnergyMetricDetails{
			{
				Name:        "total_power_milliwatts",
				Type:        "gauge",
				Description: "Total power consumption of the host in milliwatts",
				Labels:      []string{"host", "total"},
				Fn:          collectTotalEnergy,
				metricDesc:  &prometheus.Desc{},
			},
		},
		inChannelITResources: inChannelITResources,
		OutChannelEnergy:     make(chan *schema.Energy),
	}, nil
}

var implementedProviders = map[string]EnergyProvider{
	ENERGY_PROVIDER_SCAPHANDRE: scaphandre.Provider{},
}

func (g EnergyP) ListenOnChannel() {
	for r := range g.inChannelITResources {
		slog.Info("received IT resource", slog.String("name", r.Name))
	}
}

func (c EnergyProviderConfig) VerifyConfig(kind string) error {
	regexpIP, err := regexp.Compile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	if err != nil {
		return fmt.Errorf("error while compiling IP regex: %v", err)
	}

	regexpURL, err := regexp.Compile(`^(([^:\/?#]+):)?(\/\/([^\/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)
	if err != nil {
		return fmt.Errorf("error while compiling URL regex: %v", err)
	}
	var multiErr error
	for i := range c.Endpoints {
		if regexpIP.MatchString(c.Endpoints[i]) && regexpURL.MatchString(c.Endpoints[i]) {
			multiErr = multierror.Append(multiErr, fmt.Errorf("invalid endpoint: %s", c.Endpoints[i]))
		}
	}
	if multiErr != nil {
		return multiErr
	}
	// check if provider is implemented
	if _, ok := implementedProviders[kind]; !ok {
		return fmt.Errorf("provider %s is not implemented", kind)
	}
	return nil
}

func (p EnergyP) Describe(ch chan<- *prometheus.Desc) {
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

func (p EnergyP) Collect(ch chan<- prometheus.Metric) {
	slog.Info("collecting energy metrics")

	energyRecords, err := p.requestEnergy(p.config)
	if err != nil {
		slog.Error("error while collecting energy data", slog.String("error", err.Error()))
		return
	}

	// TODO: add this feature later
	// slog.Info("sending energy data to other providers")
	// for j := range energyRecords {
	// 	p.OutChannelEnergy <- energyRecords[j]
	// }

	slog.Info("register metrics")
	for j := range p.metrics {
		ch <- prometheus.MustNewConstMetric(p.metrics[0].metricDesc, prometheus.GaugeValue, float64(p.metrics[j].Fn(energyRecords)), p.metrics[j].Labels...)
	}
}

func (p EnergyP) requestEnergy(cfg EnergyProviderConfig) ([]*schema.Energy, error) {
	slog.Info("get energy", "endpoint", cfg.Endpoints)
	var multiError error
	var energyRecords []*schema.Energy
	for i := range cfg.Endpoints {
		e, err := implementedProviders[p.kind].GetEnergy(cfg.Endpoints[i])
		if err != nil {
			multiError = multierror.Append(multiError, err)
		}
		energyRecords = append(energyRecords, e...)
	}
	slog.Info("collected energy records", slog.Int("number of energy records", len(energyRecords)))
	return energyRecords, nil
}

type EnergyMetricDetails struct {
	Name        string
	Type        string
	Description string
	Labels      []string
	Fn          func(energy []*schema.Energy) float64
	metricDesc  *prometheus.Desc
}

func collectTotalEnergy(energyRecords []*schema.Energy) float64 {
	slog.Info("collecting total host energy", slog.Int("number of energy records", len(energyRecords)))
	var totalEnergyMicroWatts float64
	for i := range energyRecords {
		slog.Info("energy record", slog.String("name", energyRecords[i].Name), slog.Float64("amount", energyRecords[i].Amount), slog.String("unit", string(energyRecords[i].Unit)))
		totalEnergyMicroWatts += energyRecords[i].ConvertToMilliwatt()
	}
	return totalEnergyMicroWatts
}
