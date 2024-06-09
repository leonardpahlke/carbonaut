#!/bin/bash

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

POD_NAME=$(kubectl get pods -n carbonaut -l app=carbonaut -o jsonpath="{.items[0].metadata.name}")

if [ -z "$POD_NAME" ]; then
  echo "No pod found with label app=carbonaut"
  exit 1
fi

echo "Pod Name: $POD_NAME"
kubectl logs $POD_NAME -n carbonaut
kubectl port-forward $POD_NAME 8088:8088 -n carbonaut

# curl localhost:8088/metrics-json > dev/s2-s1-metrics.json
# curl ocalhost:8088/static-data > dev/s2-s1-state.json
