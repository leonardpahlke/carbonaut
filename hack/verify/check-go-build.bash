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

# This is the concurrency limit
MAX_POOL_SIZE=3
# This is used within the program. Do not change.
CURRENT_POOL_SIZE=0

# Jobs will be loaded from this file
PLATFORMS=(
    linux/amd64
    linux/386
    linux/arm
    linux/arm64
    linux/ppc64le
    linux/s390x
    windows/amd64
    windows/386
    freebsd/amd64
    darwin/amd64
)

for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%/*}"
    ARCH=$(basename "$PLATFORM")

    # This is the blocking loop where it makes the program to wait if the job pool is full
    while [ $CURRENT_POOL_SIZE -ge $MAX_POOL_SIZE ]; do
        CURRENT_POOL_SIZE=$(jobs | wc -l)
    done
    echo " $(date +'[%F %T]') - Build on platform $PLATFORM"
    GOARCH="$ARCH" GOOS="$OS" go build ./... &
    CURRENT_POOL_SIZE=$(jobs | wc -l)
done

# wait for all background jobs (forks) to exit before exiting the parent process
wait
