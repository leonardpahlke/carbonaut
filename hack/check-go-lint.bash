#!/bin/bash

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

VERSION=v1.52.1
URL_BASE=https://raw.githubusercontent.com/golangci/golangci-lint
URL=$URL_BASE/$VERSION/install.sh

if [[ ! -f .golangci.yaml ]]; then
    echo " $(date +'[%F %T]') - ERROR: missing .golangci.yaml in repo root" >&2
    exit 1
fi

if ! command -v golangci-lint; then
    curl -sfL $URL | sh -s $VERSION
    PATH=$PATH:bin
fi

golangci-lint version
golangci-lint linters
golangci-lint run "$@"
