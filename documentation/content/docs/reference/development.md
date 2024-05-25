---
weight: 6
---

## **Carbonaut Project Development Guide**

### Install required tools

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

#### Testing

If access keys to external sources are not set, tests will either set a mock or skip the test. To run all tests the following information needs to be set as environment variables.

- `ELECTRICITY_MAP_AUTH_TOKEN` needs to be set to run all integration tests for the electricity map provider
- `METAL_AUTH_TOKEN` needs to be set to run all integration tests for the equinix provider
