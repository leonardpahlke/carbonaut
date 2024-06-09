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
COSIGN_KEY="$HOME/.cosign/cosign.key"
COSIGN_PUB_KEY="$HOME/.cosign/cosign.pub"
BUILDER_NAME="carbonautbuilder"

if docker buildx inspect $BUILDER_NAME > /dev/null 2>&1; then
  echo "Remove existing builder"
  docker buildx rm $BUILDER_NAME
fi

echo "Create and use a new buildx builder"
docker buildx create --name $BUILDER_NAME --use
docker buildx inspect $BUILDER_NAME --bootstrap

echo "Build and push the image for multiple platforms"
docker buildx build -f Containerfile --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 -t $IMAGE_NAME:$TAG --push .

get_image_digests() {
  docker manifest inspect $IMAGE_NAME:$TAG | jq -r '.manifests[].digest'
}

echo "Waiting for the image to be available in the registry..."
for i in {1..10}; do
  IMAGE_DIGESTS=$(get_image_digests)
  if [[ -n "$IMAGE_DIGESTS" && "$IMAGE_DIGESTS" != "null" ]]; then
    break
  fi
  echo "Retrying in 5 seconds..."
  sleep 5
done

if [[ -z "$IMAGE_DIGESTS" || "$IMAGE_DIGESTS" == "null" ]]; then
  echo "Error: Image digests not found. Exiting."
  exit 1
fi

echo $IMAGE_DIGESTS

# echo "Check if Cosign key exists, if not, generate a new key pair"
# if [ ! -f "$COSIGN_KEY" ]; then
#     cosign generate-key-pair
# fi

# echo "Iterate over each digest to generate SBOM and attest the image"
# for DIGEST in $IMAGE_DIGESTS; do
#   echo "Generate the SBOM for digest: $DIGEST"
#   SBOM_FILE="sbom-$DIGEST.json"
#   syft $IMAGE_NAME@$DIGEST -o syft-json > $SBOM_FILE

#   echo "Attest the image with the SBOM for digest: $DIGEST"
#   cosign attest --key $COSIGN_KEY --predicate $SBOM_FILE --type https://spdx.dev/Document $IMAGE_NAME@$DIGEST

#   echo "Verify the attestation for digest: $DIGEST"
#   cosign verify-attestation --key $COSIGN_PUB_KEY $IMAGE_NAME@$DIGEST

#   echo "SBOM generation and attestation completed for digest: $DIGEST"
# done

# echo "Build, SBOM generation, and push completed."
