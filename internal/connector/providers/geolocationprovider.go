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
	"regexp"

	"carbonaut.cloud/pkg/providers/ipapi"
	"carbonaut.cloud/pkg/schema"
	"carbonaut.cloud/pkg/util/compareutils"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"
)

const (
	GEO_PROVIDER_IPAPI = "ipapi"
)

type GeoProviderConfig struct {
	// IP addresses to look up; needs to be valid IPv4 addresses
	IPAddresses []string `json:"ip_addresses"`
}

type GeoProvider interface {
	GetGeolocation(ip string) (*schema.Geolocation, error)
}

var implementedGeoProviders = map[string]GeoProvider{
	GEO_PROVIDER_IPAPI: ipapi.Provider{},
}

type GeoP struct {
	config                GeoProviderConfig
	kind                  string
	collectorName         string
	metrics               []GeoMetricDetails
	inChannelITResources  chan *schema.ITResource
	OutChannelGeolocation chan *schema.Geolocation
}

func NewGeoProvider(spec any, kind string, collectorName string, inChannelITResources chan *schema.ITResource) (*GeoP, error) {
	var cfg GeoProviderConfig
	specBytes, err := yaml.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling geo provider config: %w", err)
	}
	if err := yaml.Unmarshal(specBytes, &cfg); err != nil {
		return nil, fmt.Errorf("error while unmarshalling geo provider config: %w", err)
	}
	if err := cfg.VerifyConfig(kind); err != nil {
		return nil, fmt.Errorf("error while verifying the provider config: %w", err)
	}
	slog.Info("geo provider config verified", slog.String("kind", kind))
	return &GeoP{
		config:                cfg,
		kind:                  kind,
		collectorName:         collectorName,
		metrics:               []GeoMetricDetails{{Name: "geo_region_span", Type: "gauge", Description: "The total number of regions spanned by the IT Infrastructure", Labels: []string{"total"}, Fn: collectDistinctGeolocations, metricDesc: &prometheus.Desc{}}},
		inChannelITResources:  inChannelITResources,
		OutChannelGeolocation: make(chan *schema.Geolocation),
	}, nil
}

func (g GeoP) ListenOnChannel() {
	for r := range g.inChannelITResources {
		slog.Info("received IT resource", slog.String("name", r.Name))
	}
}

func (g GeoProviderConfig) VerifyConfig(kind string) error {
	// TODO: check if the config is empty
	r, err := regexp.Compile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
	if err != nil {
		return fmt.Errorf("error while compiling regex: %v", err)
	}
	for _, ip := range g.IPAddresses {
		if match := r.MatchString(ip); !match {
			return fmt.Errorf("invalid ip: %s", ip)
		}
	}
	// check if provider is implemented
	if _, ok := implementedGeoProviders[kind]; !ok {
		return fmt.Errorf("provider %s is not implemented", kind)
	}
	return nil
}

func (p GeoP) Describe(ch chan<- *prometheus.Desc) {
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

func (p GeoP) Collect(ch chan<- prometheus.Metric) {
	slog.Info("collecting geolocation metrics")
	geolocations, err := p.requestGeolocation(p.config)
	if err != nil {
		slog.Error("error while getting geolocation", slog.String("error", err.Error()))
		return
	}
	slog.Info("collected geolocations", slog.Int("number of geolocations", len(geolocations)))
	// TODO: add this feature later
	// slog.Info("sending geolocations to other providers")
	// for j := range geolocations {
	// 	p.OutChannelGeolocation <- geolocations[j]
	// }
	slog.Info("register metrics")
	for j := range p.metrics {
		ch <- prometheus.MustNewConstMetric(p.metrics[0].metricDesc, prometheus.GaugeValue, float64(p.metrics[j].Fn(geolocations)), p.metrics[j].Labels...)
	}
}

func (p GeoP) requestGeolocation(cfg GeoProviderConfig) ([]*schema.Geolocation, error) {
	slog.Info("get geolocation", "provider", p.kind)
	provider := implementedGeoProviders[p.kind]
	var multiError error
	var geolocations []*schema.Geolocation
	for i := range cfg.IPAddresses {
		geolocation, err := provider.GetGeolocation(cfg.IPAddresses[i])
		if err != nil {
			multiError = multierror.Append(multiError, err)
		}
		geolocations = append(geolocations, geolocation)
	}
	if multiError != nil {
		return geolocations, fmt.Errorf("error while getting geolocation: %w", multiError)
	}
	slog.Info("collected geolocations", slog.Any("geolocations", geolocations[0]))
	return geolocations, nil
}

type GeoMetricDetails struct {
	Name        string
	Type        string
	Description string
	Labels      []string
	Fn          func(geolocations []*schema.Geolocation) int
	metricDesc  *prometheus.Desc
}

func collectDistinctGeolocations(geolocations []*schema.Geolocation) int {
	slog.Info("collecting distinct geolocations")
	var distinctGeolocations []*schema.Geolocation
	for _, geolocation := range geolocations {
		if !compareutils.CheckListContains(distinctGeolocations, geolocation) {
			distinctGeolocations = append(distinctGeolocations, geolocation)
		}
	}
	return len(distinctGeolocations)
}
