![carbonaut-banner](.github/carbonaut-banner.png)

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/leonardpahlke/carbonaut.svg)](https://github.com/leonardpahlke/carbonaut)
[![Go Report Card](https://goreportcard.com/badge/leonardpahlke/carbonaut)](https://goreportcard.com/report/leonardpahlke/carbonaut)

[Carbonaut](https://carbonaut.dev/) is an open-source cloud native software project which aims to establish transparency for energy and IT-Resources used, emissions caused and eventually, which natural resources are used to run your software.
The project is at the start to fulfill this vision. 

## Development

**Install**:
1. Install [`pre-commit`](https://pre-commit.com/) which is used to run code checks
2. Run `make install` which installs Go packages and sets up pre-commit. 

```
$ make help
Available commands:
  all                    - Build project resources and verify code
  verify                 - Run verifications on the project (lint, vet, tests)
  install                - Install project dependencies
  format                 - Format Go files
  upgrade                - Upgrade project dependencies
  compile-grpc           - Compile gRPC and protobuf definitions
  test-coverage          - Generate and open test coverage report
  clean-coverage         - Clean test coverage reports
  tf-init                - Initialize OpenTofu configuration
  ...
```


## Implemented Providers

### Static Resource Providers:

* [Equinix](https://www.equinix.com/)

### Dynamic Resource Providers:

* [Scaphandre](https://github.com/hubblo-org/scaphandre)

### Dynamic Environment Providers:

* [Electricity Map](https://www.electricitymaps.com): (you can get a free tier account [here](https://www.electricitymaps.com/pricing))


## Testing

If access keys to external sources are not set, tests will either set a mock or skip the test. To run all tests the following information needs to be set as environment variables.
* `ELECTRICITY_MAP_AUTH_TOKEN` needs to be set to run all integration tests for the electricity map provider
* `METAL_AUTH_TOKEN` needs to be set to run all integration tests for the equinix provider


## TODO

Build own stress test container images since polinux/stress and yauritux/sysbench are not maintained. Smth like this could work.

```Dockerfile
FROM debian:12

RUN apt-get update && apt-get install -y stress \
        --no-install-recommends && rm -r /var/lib/apt/lists/*

CMD ["stress", "--verbose", "--vm", "1", "--vm-bytes", "256M"]
```

---

Extend Equinix provider to support paging in resource & project discovery. Paging information is provided in the "Meta" information which right now is not parsed and processed.

---

The energy mix data could be cached across resources. It's likely that multiple resources are deployed in the same region - its therefore not needed to query data by resource.
