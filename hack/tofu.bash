#!/bin/bash

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

USE_DEFAULTS="${USE_DEFAULTS:-0}"
SSH_KEY_PATH="${SSH_KEY_PATH:-$HOME/.ssh/id_equinix_carbonaut_ed25519.pub}"
PRIVATE_KEY_PATH="${PRIVATE_KEY_PATH:-$HOME/.ssh/id_equinix_carbonaut_ed25519}"
CARBONAUT_NUM_PROJECTS="${CARBONAUT_NUM_PROJECTS:-1}"
CARBONAUT_VM_COUNT_PROJECTS="${CARBONAUT_VM_COUNT_PROJECTS:-1}"

function ask_ssh_key() {
  if [ "$USE_DEFAULTS" -eq 0 ]; then
    echo "Current SSH key path: $SSH_KEY_PATH"
    read -p "Enter new SSH key path or press enter to use default: " input_key
    if [ "$input_key" != "" ]; then
      SSH_KEY_PATH=$input_key
    fi
  else
    echo "Using default SSH key path: $SSH_KEY_PATH"
  fi
}

function ask_private_key() {
  if [ "$USE_DEFAULTS" -eq 0 ]; then
    echo "Current private key path: $PRIVATE_KEY_PATH"
    read -p "Enter new private key path or press enter to use default: " input_key
    if [ "$input_key" != "" ]; then
      PRIVATE_KEY_PATH=$input_key
    fi
  else
    echo "Using default private key path: $PRIVATE_KEY_PATH"
  fi
}

function tf_vars() {
  if [ "$USE_DEFAULTS" -eq 0 ]; then
    echo "Current num_projects: $CARBONAUT_NUM_PROJECTS"
    read -p "Enter how many projects you like to deploy: " input_projects
    if [ "$input_projects" != "" ]; then
      CARBONAUT_NUM_PROJECTS=$input_projects
    fi
  fi

  if [ "$USE_DEFAULTS" -eq 0 ]; then
    echo "Current vm_count_per_project: $CARBONAUT_VM_COUNT_PROJECTS"
    read -p "Enter how many resources you like to deploy for each project: " input_vms
    if [ "$input_vms" != "" ]; then
      CARBONAUT_VM_COUNT_PROJECTS=$input_vms
    fi
  fi
}

function show_help() {
  echo "Usage: $0 {configure|apply|destroy|plan|connect|stress-test}"
  echo ""
  echo "Commands:"
  echo "  plan          Plan the Terraform configuration"
  echo "  apply         Apply the Terraform configuration"
  echo "  destroy       Destroy the Terraform configuration"
  echo "  configure     Run the Ansible setup for all created resources"
  echo "  connect       Connect to a VM using SSH"
  echo "  stress-test   Run stress test script on all configured machines"
  exit 0
}

if [ $# -eq 0 ]; then
  show_help
fi

case $1 in
  configure)
    ask_private_key
    OUTPUT=$(tofu -chdir=test-scenario/dev output -json vm_public_ip)
    IP_ADDRESSES=$(echo $OUTPUT | jq -r '.[]')
    for IP in $IP_ADDRESSES; do
      ansible-playbook -i "$IP," test-scenario/dev/setup_vm.yml -u root --private-key="$PRIVATE_KEY_PATH"
    done
    ;;
  apply)
    ask_ssh_key
    tf_vars
    tofu -chdir=test-scenario/dev apply -var "public_key=$(cat $SSH_KEY_PATH)" -var "num_projects=$CARBONAUT_NUM_PROJECTS" -var "vm_count_per_project=$CARBONAUT_VM_COUNT_PROJECTS"
    ;;
  destroy)
    ask_ssh_key
    tofu -chdir=test-scenario/dev destroy -var "public_key=$(cat $SSH_KEY_PATH)" -var "num_projects=$CARBONAUT_NUM_PROJECTS" -var "vm_count_per_project=$CARBONAUT_VM_COUNT_PROJECTS"
    ;;
  plan)
    ask_ssh_key
    tf_vars
    tofu -chdir=test-scenario/dev plan -var "public_key=$(cat $SSH_KEY_PATH)" -var "num_projects=$CARBONAUT_NUM_PROJECTS" -var "vm_count_per_project=$CARBONAUT_VM_COUNT_PROJECTS"
    ;;
  connect)
    ask_private_key
    OUTPUT=$(tofu -chdir=test-scenario/dev output -json vm_public_ip)
    IP_ADDRESSES=($(echo $OUTPUT | jq -r '.[]'))
    
    echo "Select the IP address to connect to:"
    for i in "${!IP_ADDRESSES[@]}"; do
      echo "$i: ${IP_ADDRESSES[$i]}"
    done
    
    read -p "Enter the index of the IP address you want to connect to: " index
    SERVER_IP=${IP_ADDRESSES[$index]}
    
    ssh -i "$PRIVATE_KEY_PATH" root@"$SERVER_IP"
    ;;
  connection-verify)
    ask_private_key
    OUTPUT=$(tofu -chdir=test-scenario/dev output -json vm_public_ip)
    IP_ADDRESSES=($(echo $OUTPUT | jq -r '.[]'))
    
    for i in "${!IP_ADDRESSES[@]}"; do
      if curl --output /dev/null --silent --head --fail "${IP_ADDRESSES[$i]}:8080"; then
        echo "✅ ${IP_ADDRESSES[$i]}"
      else
        echo "❌ ${IP_ADDRESSES[$i]}"
      fi
    done
    ;;
  stress-test)
    ask_private_key
    OUTPUT=$(tofu -chdir=test-scenario/dev output -json vm_public_ip)
    IP_ADDRESSES=$(echo $OUTPUT | jq -r '.[]')
    for IP in $IP_ADDRESSES; do
      ansible-playbook -i "$IP," test-scenario/dev/stress_test.yml -u root --private-key="$PRIVATE_KEY_PATH"
    done
    ;;
  *)
    show_help
    ;;
esac
