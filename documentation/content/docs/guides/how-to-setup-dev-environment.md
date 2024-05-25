---
weight: 2
---

## **Carbonaut: How To Setup Your Dev Environment**

This guide gives information how to start setup your development environment so you are ready to customize the current Carbonaut version.

There are two options which you can use. Either you use a vscode dev container which uses Docker in the background as virtualized dev environment. Or you use your regular machine for development. If you are on Windows you need to work with [WSL (Linux Subsystem for Windows)](https://learn.microsoft.com/en-us/windows/wsl/install) and note that there may be some issues since Carbonaut was developed on macOS and not tested on other platforms. If any step is not working, please open a PR to improve this document.

**Make sure to work on a fork and not the cloned carbonaut repository! See [contributor](/docs/reference/contributing/) guide for more information.**

### 1. Manual Setup

1. If you are on Windows install [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) and work going forward on the Linux subsystem.
2. Install [Go](https://go.dev/doc/install) and [NPM](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm). Go is used to compile the project and install additional tooling like [golangci-lint](https://github.com/golangci/golangci-lint). NPM is used to install additional tooling like [pre-commit](https://pre-commit.com/) and [prettier](https://prettier.io/) for formatting. All installs are listed in the `Makefile` under `installs`
   1. Current Go version used on macOS `go version go1.22.2 darwin/arm64`
   2. Current NPM version used on macOS `10.8.0`
3. Install dependencies with `make install`.

### 2. Use VSCode Dev Containers

{{< hint info >}}
[**Dev Containers**](https://code.visualstudio.com/docs/devcontainers/containers) are not yet implemented.
{{< /hint >}}
