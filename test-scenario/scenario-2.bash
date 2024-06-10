#!/bin/bash

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

SSH_KEY_PATH="${SSH_KEY_PATH:-$HOME/.ssh/id_equinix_carbonaut_ed25519.pub}"
CARBONAUT_DEFAULT_PORT="${CARBONAUT_DEFAULT_PORT:-8088}"
CARBONAUT_DEFAULT_IP="${CARBONAUT_DEFAULT_IP:-127.0.0.1}"
RESULTS_FOLDER="${RESULTS_FOLDER:-./test-scenario/results-2}"
PRIVATE_KEY_PATH="${PRIVATE_KEY_PATH:-$HOME/.ssh/id_equinix_carbonaut_ed25519}"

ENDPOINT_STATE="/static-data"
ENDPOINT_STOP="/stop"
ENDPOINT_JSON_METRICS="/metrics-json"

mkdir -p $RESULTS_FOLDER

function update_known_hosts {
    OUTPUT=$(tofu -chdir=test-scenario/dev output -json vm_public_ip)
    IP_ADDRESSES=$(echo $OUTPUT | jq -r '.[]')
    for ip in $IP_ADDRESSES; do
        echo "[s2-sx-x] updating known_hosts SSH for $ip"
        ssh-keygen -R $ip
        ssh-keyscan -H $ip >> ~/.ssh/known_hosts
    done
}

function collect_carbonaut_state_and_metrics {
    local state_file_nr=$1
    local metrics_file_nr=$2

    sleep 5
    
    echo "[$state_file_nr] collecting carbonaut state"
    curl http://$CARBONAUT_DEFAULT_IP:$CARBONAUT_DEFAULT_PORT$ENDPOINT_STATE > $RESULTS_FOLDER/$state_file_nr-state.json
    wait

    echo "[$metrics_file_nr] collecting carbonaut metrics json"
    curl http://$CARBONAUT_DEFAULT_IP:$CARBONAUT_DEFAULT_PORT$ENDPOINT_JSON_METRICS > $RESULTS_FOLDER/$metrics_file_nr-metrics.json
    wait
}

function configure_and_verify_resources {
    local configure_nr=$1
    local verify_nr=$2
    local benchmark_nr=$3

    echo "[$configure_nr] configure created resource"
    printf "$SSH_KEY_PATH\n" | make tf-configure > $RESULTS_FOLDER/$configure_nr-configure.txt
    wait

    sleep 5

    echo "[$verify_nr] verify that the resource got configured"
    printf "$SSH_KEY_PATH\n" | make tf-connection-verify > $RESULTS_FOLDER/$verify_nr-connection-verify.txt
    wait

    echo "[$benchmark_nr] start stress configured resources"
    printf "$SSH_KEY_PATH\n" | make tf-stress-test > $RESULTS_FOLDER/$benchmark_nr-stress-resources.txt
    wait
}

#### #### #### #### ####
#### SZENARIO START

echo "STARTING SZENATIO 2"

#### #### #
#### STEP 1
echo ""
echo "[s2-s1-0] STEP 1 - starting initialization step"
echo "[s2-s1-1] Setup Kubernetes cluster"
minikube start > $RESULTS_FOLDER/s2-s1-1-kubernetes-up.txt

echo "[s2-s1-2] Create Equinix Infrastructure 1xProject 1xResource"
printf "$SSH_KEY_PATH\n1\n1\nyes\n" | make tf-apply > $RESULTS_FOLDER/s2-s1-2-tofu-apply.txt
wait

update_known_hosts
wait

configure_and_verify_resources "s2-s1-3" "s2-s1-4" "s2-s1-5"

echo "[s2-s1-6] configure kubernetes cluster with carbonaut resources"
./test-scenario/deploy-k8s.bash > "$RESULTS_FOLDER/s2-s1-6-kubernetes-configure.txt"

sleep 30

POD_NAME=$(kubectl get pods -n carbonaut -l app=carbonaut -o jsonpath="{.items[0].metadata.name}")

echo "[s2-s1-7] port-forward to access carbonaut pod"
osascript <<EOF
tell application "Terminal"
    do script "cd '$(pwd)' && kubectl port-forward $POD_NAME 8088:8088 -n carbonaut > $RESULTS_FOLDER/s2-s1-7-kubernetes-port-forward.txt 2>&1; exit"
end tell
EOF

sleep 5

collect_carbonaut_state_and_metrics "s2-s1-8" "s2-s1-9"

#### #### #### #### ####
#### SZENARIO CLEAN UP
echo ""
echo "[s2-s2-0] CLEAN UP"
echo "[s2-s2-1] Remove created kubernetes resources"
kubectl delete all --all -n carbonaut

echo "[s2-s2-2] Shut down infrastructure"
printf "$SSH_KEY_PATH\nyes\n"  | make tf-destroy > $RESULTS_FOLDER/s2-s3-2-tofu-destroy.txt
wait

echo "[s2-s2-3] Stop the kubernetes cluster"
minikube stop
