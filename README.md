![carbonaut-banner](.github/carbonaut-banner.png)

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/leonardpahlke/carbonaut.svg)](https://github.com/leonardpahlke/carbonaut)
[![Go Report Card](https://goreportcard.com/badge/leonardpahlke/carbonaut)](https://goreportcard.com/report/leonardpahlke/carbonaut)

[Carbonaut](https://carbonaut.dev/) is an open source cloud native software project which aims to establish transparency for energy and IT-Resources used, emissions caused and eventually, which natural resources are used to run your software.
The project is a POC published end of May 2024 and may not get developed further.

## Information about the project is available on the carbonaut.dev website

* **ARCHITECTURE**: [link](https://carbonaut.dev/docs/concepts/components/).
* **API SCHEMA**: [link](https://carbonaut.dev/docs/reference/server-api/).
* **DATA SCHEMA**: [link](https://carbonaut.dev/docs/reference/schema/).
* **DEVELOPMENT AND CONTRIBUTION**: [link](https://carbonaut.dev/docs/reference/contributing/).
* **INSTALLATION AND DEPLOYMENT**: [link](docs/installation/getting-started/)

## OPEN ISSUES / TODO

### (S) Own Stress Test Container images

 Build own stress test container images since polinux/stress and yauritux/sysbench are not maintained. Smth like this could work.

```Dockerfile
FROM debian:12

RUN apt-get update && apt-get install -y stress \
        --no-install-recommends && rm -r /var/lib/apt/lists/*

CMD ["stress", "--verbose", "--vm", "1", "--vm-bytes", "256M"]
```


### (S) Provider Equinix Paging

Extend Equinix provider to support paging in resource & project discovery. Paging information is provided in the "Meta" information which right now is not parsed and processed.

### (S) Cache energy mix data

The energy mix data could be cached across resources. It's likely that multiple resources are deployed in the same region - its therefore not needed to query data by resource.

### (M) Detect static resource updates of remaining resources

Some of the static data may be updated without changing the resource name. This is not picked up by carbonaut right now.

### (M) Implement prometheus exporter

Carbonaut just exposes metrics in json format. This should get extended in the future to a [prometheus exporter](https://prometheus.io/docs/concepts/metric_types/).
