.PHONY: all verify build install upgrade clean compile-grpc test-coverage clean-coverage tf-init tf-plan tf-apply tf-destroy tf-connect

# Default target executed
all: build

# Verify the project code (linting, testing, checking git state)
verify:
	@echo "Verifying the project code..."
	@./hack/verify/check-git
	@go vet ./...
	@go mod tidy
	@./hack/verify/check-go-build.bash
	@./hack/verify/check-go-lint.bash
	@./hack/verify/check-go-test.bash

# Build project resources
build: compile-grpc
	@echo "Building project resources..."

# Install project dependencies
install:
	@echo "Installing project dependencies..."
	@go get ./...
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Upgrade project dependencies
upgrade:
	@echo "Upgrading project dependencies..."
	@go mod tidy
	@go get -u -t ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@go clean -i ./...

# Run test coverage report
test-coverage:
	@echo "Checking test coverage and exporting report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@open coverage.html

# Clean test coverage report
clean-coverage:
	@echo "Cleaning test coverage reports..."
	@rm -f coverage.out coverage.html

########################################
### OPEN TOFU

# Default SSH key path
SSH_KEY_PATH ?= $(HOME)/.ssh/id_equinix_carbonaut_ed25519.pub
PRIVATE_KEY_PATH ?= $(HOME)/.ssh/id_equinix_carbonaut_ed25519

# Allow user to override the default SSH key path by input during runtime
ask-ssh-key:
	@if [ -z "$${USE_DEFAULT_KEYS}" ]; then \
		echo "Current SSH key path: $(SSH_KEY_PATH)"; \
		read -p "Enter new SSH key path or press enter to use default: " input_key; \
		if [ "$$input_key" != "" ]; then \
			export SSH_KEY_PATH=$$input_key; \
		fi; \
	else \
		echo "Using default SSH key path: $(SSH_KEY_PATH)"; \
	fi;

# Allow user to override the default private key path by input during runtime
ask-private-key:
	@if [ -z "$${USE_DEFAULT_KEYS}" ]; then \
		echo "Current private key path: $(PRIVATE_KEY_PATH)"; \
		read -p "Enter new private key path or press enter to use default: " input_key; \
		if [ "$$input_key" != "" ]; then \
			export PRIVATE_KEY_PATH=$$input_key; \
		fi; \
	else \
		echo "Using default private key path: $(PRIVATE_KEY_PATH)"; \
	fi;

# OpenTofu initialization
tf-init:
	tofu -chdir=dev init

# OpenTofu planning
tf-plan: ask-ssh-key
	tofu -chdir=dev plan -var "public_key=$$(cat $(SSH_KEY_PATH))"

# OpenTofu apply
tf-apply: ask-ssh-key
	tofu -chdir=dev apply -var "public_key=$$(cat $(SSH_KEY_PATH))"

# OpenTofu Destroy
tf-destroy: ask-ssh-key
	tofu -chdir=dev destroy -var "public_key=$$(cat $(SSH_KEY_PATH))"

# Fetch the IP address from OpenTofu and connect
tf-connect: ask-private-key
	$(eval SERVER_IP := $(shell tofu -chdir=dev output -raw device_public_ip))
	ssh -i $(PRIVATE_KEY_PATH) root@$(SERVER_IP)

ansible-setup: ask-private-key
	$(eval SERVER_IP := $(shell tofu -chdir=dev output -raw device_public_ip))
	ansible-playbook -i $(SERVER_IP), dev/setup_vm.yml -u root --private-key=$(PRIVATE_KEY_PATH)

########################################
### GENERAL

help:
	@echo "Available commands:"
	@echo "  all                    - Build project resources and verify code"
	@echo "  verify                 - Run verifications on the project (lint, vet, tests)"
	@echo "  build                  - Build project resources"
	@echo "  install                - Install project dependencies"
	@echo "  upgrade                - Upgrade project dependencies"
	@echo "  clean                  - Clean build artifacts and dependencies"
	@echo "  compile-grpc           - Compile gRPC and protobuf definitions"
	@echo "  test-coverage          - Generate and open test coverage report"
	@echo "  clean-coverage         - Clean test coverage reports"
	@echo "  tf-init                - Initialize OpenTofu configuration"
	@echo "  tf-plan                - Plan OpenTofu changes"
	@echo "  tf-apply               - Apply OpenTofu changes"
	@echo "  tf-destroy             - Destroy the created OpenTofu infrastrucutre"
	@echo "  tf-connect             - Connect to the created server"
	@echo "  SSH_KEY_PATH           - Current SSH key path: $(SSH_KEY_PATH)"
	@echo "  ansible-setup          - Setup the machine with required packages etc."
	@echo "  PRIVATE_KEY_PATH       - Current private SSH key path: $(PRIVATE_KEY_PATH)"
