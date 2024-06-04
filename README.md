![carbonaut-banner](https://carbonaut.dev/carbonaut-banner.png)

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/leonardpahlke/carbonaut.svg)](https://github.com/leonardpahlke/carbonaut)
[![Go Report Card](https://goreportcard.com/badge/leonardpahlke/carbonaut)](https://goreportcard.com/report/leonardpahlke/carbonaut)

[Carbonaut](https://carbonaut.dev/) is project to collect and refine environmental sustainability data from your IT infrastructure and make it available in a common data schema. The project does not implement a scraper to scan your IT infrastructure or measure your virtual machine energy, but integrates over providers which are implemented as plugins to integrate external data sources. The project is a POC and therefore not production ready.

[DOCS](https://carbonaut.dev/docs/) - [ARCHITECTURE](https://carbonaut.dev/docs/components/)

## OPEN ISSUES / TODO

* **(S) Provider Equinix Paging**: Extend Equinix provider to support paging in resource & project discovery. Paging information is provided in the "Meta" information which right now is not parsed and processed.
* **(S) Cache energy mix data**: The energy mix data could be cached across resources. It's likely that multiple resources are deployed in the same region - its therefore not needed to query data by resource.
* **(M) Detect static resource updates of remaining resources**: Some of the static data may be updated without changing the resource name. This is not picked up by carbonaut right now.
* **(M) Implement prometheus exporter**: Carbonaut just exposes metrics in json format. This should get extended in the future to a [prometheus exporter](https://prometheus.io/docs/concepts/metric_types/).
