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

package scaphandre

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"carbonaut.cloud/pkg/schema"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// Scaphandre is an energy provider that scrapes prometheus metrics from a Scaphandre server

var scrapeKeys []string = []string{
	"scaph_host_power_microwatts",
}

type Provider struct{}

func (Provider) GetEnergy(endpoint string) ([]*schema.Energy, error) {
	slog.Debug("prepare http request")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while scraping metrics: %w", err)
	}
	defer resp.Body.Close()

	energyMetricsMap, err := parsePromMetricStr(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while parsing metrics: %w", err)
	}
	slog.Info("scraped metrics", slog.Int("count", len(energyMetricsMap)))

	energyRecords := []*schema.Energy{}
	for metricName := range energyMetricsMap {
		energyRecords = append(energyRecords, &schema.Energy{
			Amount: energyMetricsMap[metricName].GetMetric()[0].GetGauge().GetValue(),
			Unit:   schema.MICROWATT,
			Name:   metricName,
		})
	}
	return energyRecords, nil
}

// this function parses a string of prometheus metrics and returns a map of metrics
// example read in string:
// # HELP scaph_host_power_microwatts Total power consumption of the host in microwatts
// # TYPE scaph_host_power_microwatts gauge
// scaph_host_power_microwatts{host="carbonaut"} 123456789
func parsePromMetricStr(r io.Reader) (map[string]*dto.MetricFamily, error) {
	scanner := bufio.NewScanner(r)
	helpLine := ""
	typeLine := ""
	filteredMetrics := ""
	for scanner.Scan() {
		currentLine := scanner.Text()
		for i := range scrapeKeys {
			if strings.HasPrefix(currentLine, scrapeKeys[i]) {
				filteredMetrics = filteredMetrics + "\n" + helpLine + "\n" + typeLine + "\n" + currentLine + "\n"
			}
		}
		helpLine = typeLine
		typeLine = currentLine
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning metrics: %w", err)
	}

	var parser expfmt.TextParser
	scrapedMetrics, err := parser.TextToMetricFamilies(strings.NewReader(filteredMetrics))
	if err != nil {
		return nil, fmt.Errorf("error while parsing metrics: %w", err)
	}
	return scrapedMetrics, nil
}
