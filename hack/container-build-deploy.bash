#!/bin/bash

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

IMAGE_NAME="leonardpahlke/carbonaut"
TAG="latest"
SBOM_FILE="sbom.json"

docker buildx create --use -f Containerfile
docker buildx inspect --bootstrap -f Containerfile
docker buildx build -f Containerfile --platform linux/amd64,linux/arm64 -t $IMAGE_NAME:$TAG --push .

IMAGE_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' $IMAGE_NAME:$TAG)

syft $IMAGE_DIGEST -o syft-json > $SBOM_FILE

# if [ ! -f "$COSIGN_KEY" ]; then
#     cosign generate-key-pair
# fi

cosign attest --key $COSIGN_KEY --predicate $SBOM_FILE --type https://spdx.dev/Document $IMAGE_DIGEST

cosign verify-attestation --key cosign.pub $IMAGE_DIGEST

echo "Build, SBOM generation, and push completed."
