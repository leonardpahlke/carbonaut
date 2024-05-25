---
weight: 5
---

## **Carbonaut Contributing**

The Carbonaut project is a POC project and therefore no community structures are in place. If you find this project interesting enough to contribute, please open up an issue on the repository or directly a PR to discuss your idea.
Any contributions are very welcome!

### Development Workflow

* Fork the repository and work on your fork. It's recommended to create a feature branch on your fork and open pull requests from feature branches to Carbonaut's main branch. If you have questions about forks, branches etc. take a look at [GitHub's documentation](https://docs.github.com/en).
* If you forked the repository, install all dependencies, Go, `pre-commit` and then run `make install` to install other go nbased tools (see [`Makefile`](https://github.com/leonardpahlke/carbonaut/blob/main/Makefile)). If you intend to make changes to the manual testing scenario, refer to this [guide](https://github.com/leonardpahlke/carbonaut/blob/main/dev/README.md).
* After that you can run `make verify` to check if everything is setup for local development.

### Getting setup for development


#### Install required tools
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
* `ELECTRICITY_MAP_AUTH_TOKEN` needs to be set to run all integration tests for the electricity map provider
* `METAL_AUTH_TOKEN` needs to be set to run all integration tests for the equinix provider


### Ways of communication

* GitHub issues and pull requests on the Carbonaut repository (no forks!)

### Additional Comments by area of contribution

* **Improve internal code**: improvments are welcome! There are several `TODO: XYZ` annotations in the code that highlight some areas of improvements.
* **Increasing test coverage and quality**: improvements are welcome! The test coverage is not great. Test coverage is uploaded as artifact with each push to the main branch (see actions). You can also run test coverage of the Go code by executing `make test-coverage`. 
* **Adding Provider Plugins**: make sure to test it both mocked and as E2E test.
* **Proposing changing data schema or API**: sure, make your case. The datamodels are minimal and not complete, changes are welcome.

Improvements to this document are welcome!
