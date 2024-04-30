#!/bin/bash

# Copyright 2023 CARBONAUT AUTHORS
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
