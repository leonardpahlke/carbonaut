#!/bin/bash

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

NAMESPACE="carbonaut"

if ! kubectl get namespace $NAMESPACE; then
  echo "Creating namespace $NAMESPACE"
  kubectl create namespace $NAMESPACE
else
  echo "Namespace $NAMESPACE already exists. Skipping creation."
fi

echo "Create or update the secret"
kubectl delete secret carbonaut-secrets -n $NAMESPACE --ignore-not-found

kubectl create secret generic carbonaut-secrets \
  --from-literal=METAL_AUTH_TOKEN=$METAL_AUTH_TOKEN \
  --from-literal=ELECTRICITY_MAP_AUTH_TOKEN=$ELECTRICITY_MAP_AUTH_TOKEN \
  -n $NAMESPACE

echo "Apply the Kubernetes configuration"
kubectl apply -f dev/k8s.yaml -n $NAMESPACE
