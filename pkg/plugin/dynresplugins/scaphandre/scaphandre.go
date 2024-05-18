package scaphandre

import (
	"errors"
	"fmt"
	"strconv"

	"carbonaut.dev/pkg/plugin"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/types/dynres"
	"carbonaut.dev/pkg/util/promscraper"
)

var PluginName plugin.Kind = "scaphandre"

type p struct {
	cfg *dynres.Config
}

func New(cfg *dynres.Config) (p, error) {
	return p{
		cfg: cfg,
	}, nil
}

func (p) GetName() *plugin.Kind {
	return &PluginName
}

func (p p) GetDynamicResourceData(data *resource.StaticResData) (*resource.DynamicResData, error) {
	if data.IPv4 == "" {
		return nil, errors.New("IPv4 is empty and therefore no call can be made to collect scaphandre data")
	}
	metricData, err := collectSpecifiedMetricValues(promscraper.Prom2Json{
		URL:               "http://" + data.IPv4 + p.cfg.Endpoint,
		Cert:              p.cfg.Cert,
		Key:               p.cfg.Key,
		AcceptInvalidCert: p.cfg.AcceptInvalidCert,
	}, []string{
		"scaph_host_cpu_frequency",
		"scaph_host_energy_microjoules",
		"scaph_process_cpu_usage_percentage",
	})
	if err != nil {
		return nil, fmt.Errorf("error collecting data from scraphandre: %v", err)
	}

	return &resource.DynamicResData{
		CPUFrequency:          toFloatOrZero(metricData["scaph_host_cpu_frequency"]),
		EnergyHostMicrojoules: toIntOrZero(metricData["scaph_host_energy_microjoules"]),
		CPULoadPercentage:     toFloatOrZero(metricData["scaph_process_cpu_usage_percentage"]),
	}, nil
}

func collectSpecifiedMetricValues(cfg promscraper.Prom2Json, keys []string) (map[string]string, error) {
	familyData, err := promscraper.Collect(cfg)
	if err != nil {
		return nil, err
	}
	filteredData := make(map[string]string)
	for i := range keys {
		for j := range familyData {
			if keys[i] == familyData[j].Name {
				val := ""
				if len(familyData[j].Metrics) > 0 {
					if metric, ok := familyData[j].Metrics[0].(promscraper.Metric); ok {
						fmt.Printf("First Metric: %+v\n", metric)
						val = metric.Value
					} else {
						fmt.Println("First element is not a Metric")
					}
				} else {
					fmt.Println("No metrics available")
				}
				filteredData[keys[i]] = val
				break
			}
		}
	}
	return filteredData, nil
}

func toIntOrZero(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

func toFloatOrZero(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0
	}
	return num
}
