---
title: Develop Internals
weight: 3
slug: how-to-develop-internals
---

This guide gives information how to start develop internals of Carbonaut. Carbonauts source code is developed in Go and located in the `pkg/` directory of the project. After this guide you are able to make changes, verify if made changes work and how to push these changes and get towards merging it to the main Carbonaut project.

{{% steps %}}

### Setup

1. Read the contributors guide [here](/docs/contributing/) and fork the [carbonaut repository](https://github.com/leonardpahlke/carbonaut/fork).
2. Setup your dev environment. Follow this [guide](/docs/guides/how-to-setup-dev-environment/).
3. Test building the project locally `make build`
4. Run Carbonaut in test-run mode `./main --test-run`. Carbonaut will load a test config and shut down after a couple of seconds automatically.

### Make changes to the codebase

Open up any edit and make changes to the code base. You can run unit tests and linting with `make verify`. During that process a `coverage.html` will be generated which can be used to explore test coverage.

### Verify changes made

1. Run `make verify`.
   1. INFO: Some steps run commands like `go mod tidy` which may change `go.mod` (depending on your edits) if `go.mod` was changed during `make verify`, the check will fail but its safe to rerun it.
2. Optionally you can run e2e checks see this [guide](/docs/guides/how-to-test-e2e/)

{{% /steps %}}
