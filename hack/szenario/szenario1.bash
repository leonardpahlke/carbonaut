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
RESULTS_FOLDER="${RESULTS_FOLDER:-./dev/scenario/results-1}"
PRIVATE_KEY_PATH="${PRIVATE_KEY_PATH:-$HOME/.ssh/id_equinix_carbonaut_ed25519}"

ENDPOINT_STATE="/static-data"
ENDPOINT_STOP="/stop"
ENDPOINT_JSON_METRICS="/metrics-json"

mkdir -p $RESULTS_FOLDER

function update_known_hosts {
    OUTPUT=$(tofu -chdir=dev output -json vm_public_ip)
    IP_ADDRESSES=$(echo $OUTPUT | jq -r '.[]')
    for ip in $IP_ADDRESSES; do
        echo "[s1-sx-x] updating known_hosts SSH for $ip"
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

    echo "[$verify_nr] verify that the resource got configured"
    printf "$SSH_KEY_PATH\n" | make tf-connection-verify > $RESULTS_FOLDER/$verify_nr-configure.txt
    wait

    echo "[$benchmark_nr] start stress configured resources"
    printf "$SSH_KEY_PATH\n" | make tf-stress-test > $RESULTS_FOLDER/$benchmark_nr-stress-resources.txt
    wait
}

echo "STARTING SZENATIO 1"

#### #### #### #### ####
#### STEP 1
echo ""
echo "[s1-s1-0] STEP 1 - starting"
echo "[s1-s1-1] Create Equinix Infrastructure 1xProject 1xResource"
printf "$SSH_KEY_PATH\n1\n1\nyes\n" | make tf-apply > $RESULTS_FOLDER/s1-s1-1-carbonaut-log.txt
wait

update_known_hosts
wait

configure_and_verify_resources "s1-s1-2" "s1-s1-3" "s1-s1-4"

echo "[s1-s1-5] starting carbonaut"
osascript <<EOF
tell application "Terminal"
    do script "cd '$(pwd)' && go run main.go -c dev/config.yaml > $RESULTS_FOLDER/s1-s1-5-carbonaut-log.txt 2>&1; exit"
end tell
EOF

sleep 5

collect_carbonaut_state_and_metrics "s1-s1-6" "s1-s1-7"

# #### #### #### #### ####
# #### STEP 2
# echo ""

# echo "[S2-S2-1] dereference infrastructure in carbonaut by supplying an empty configuration"
# curl -X POST -H "Content-Type: application/x-yaml" --data-binary @dev/empty-config.yaml http://$CARBONAUT_DEFAULT_IP:$CARBONAUT_DEFAULT_PORT/load-config
# wait


# echo "[S2-S2-5] collecting carbonauts state"
# printf "$CARBONAUT_DEFAULT_PORT\n$CARBONAUT_DEFAULT_IP\n" | make carbonaut-get-state > $RESULTS_FOLDER/s2-s2-5-state.json
# wait

# echo "[S2-S2-6] collecting carbonauts metrics json"
# printf "$CARBONAUT_DEFAULT_PORT\n$CARBONAUT_DEFAULT_IP\n" | make carbonaut-get-json-metrics > $RESULTS_FOLDER/s2-s2-6-metrics.json
# wait

echo "[sX] CLEAN UP"

echo "[s1-sX-1] stop carbonaut"
curl http://$CARBONAUT_DEFAULT_IP:$CARBONAUT_DEFAULT_PORT$ENDPOINT_STOP
wait

echo "[s2-sX-2] shut down infrastructure"
printf "$SSH_KEY_PATH\nyes\n"  | make tf-destroy
wait
