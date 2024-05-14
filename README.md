![carbonaut-banner](assets/carbonaut-banner.png)

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
