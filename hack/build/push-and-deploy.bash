#!/bin/bash

# Copyright (c) 2022 CARBONAUT AUTHORS
#
# Licensed under the MIT license: https://opensource.org/licenses/MIT
# Permission is granted to use, copy, modify, and redistribute the work.
# Full license information available in the project LICENSE file.

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

# Required: authentication to container registry

# Container registry used to 
REGISTRY=ghcr.io
ORG=carbonaut-cloud

# Name of the container to build and deploy
CONTAINER=$1

echo " $(date +'[%F %T]') - :: Build container image $ORG/$CONTAINER ::"
SHORTHASH="$(git rev-parse --short HEAD)"
docker build -f build/Containerfile.$CONTAINER -t $REGISTRY/$ORG/carbonaut-$CONTAINER:latest -t $REGISTRY/$ORG/carbonaut-$CONTAINER:$SHORTHASH .
echo " $(date +'[%F %T]') - :: Push container image $REGISTRY/$ORG/carbonaut-$CONTAINER:latest:$SHORTHASH to $REGISTRY ::"
docker push $REGISTRY/$ORG/carbonaut-$CONTAINER:latest
docker push $REGISTRY/$ORG/carbonaut-$CONTAINER:$SHORTHASH
