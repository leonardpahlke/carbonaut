---
title: Setup Development Environment
weight: 2
slug: how-to-setup-dev-environment
---

This guide gives information how to start setup your development environment so you are ready to customize the current Carbonaut version.

There are two options which you can use. Either you use a vscode dev container which uses Docker in the background as virtualized dev environment. Or you use your regular machine for development. If you are on Windows you need to work with [WSL (Linux Subsystem for Windows)](https://learn.microsoft.com/en-us/windows/wsl/install) and note that there may be some issues since Carbonaut was developed on macOS and not tested on other platforms. If any step is not working, please open a PR to improve this document.

{{< callout type="warning" >}}
Make sure to work on a fork and not the cloned carbonaut repository! See the [contributor](/docs/contributing/) guide for more information.
{{< /callout >}}

{{% steps %}}

### GO & NPM

Install [Go](https://go.dev/doc/install) and [NPM](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm). Go is used to compile the project and install additional tooling like [golangci-lint](https://github.com/golangci/golangci-lint). NPM is used to install additional tooling like [pre-commit](https://pre-commit.com/) and [prettier](https://prettier.io/) for formatting. All installs are listed in the `Makefile` under `installs`
1. Current Go version used on macOS `go version go1.22.2 darwin/arm64`
2. Current NPM version used on macOS `10.8.0`
### MAKE Install

Install dependencies with `make install`.

### Optional: Use VSCode Dev Containers

{{< callout type="info" >}}
   [**Dev Containers**](https://code.visualstudio.com/docs/devcontainers/containers) are not yet implemented.
{{< /callout >}}

{{% /steps %}}

