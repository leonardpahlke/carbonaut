.PHONY: all build verify format install upgrade test-coverage clean-coverage tf-init tf tf-plan tf-apply tf-destroy tf-connect tf-configure tf-stress-test container-build-local

# Default target executed
all: verify

# Build the go binary locally
build:
	@go build main.go

# Verify the project code (linting, testing, checking git state)
verify:
	@echo "Verifying the project code..."
	@pre-commit run --all-files

# Install project dependencies
install:
	@echo "Installing project dependencies..."
	@go get ./...
	@echo "Installing Go tooling..."
	@go install github.com/4meepo/tagalign/cmd/tagalign@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install golang.org/x/tools/cmd/godoc@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.1
	@go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
	@echo "Installing Node tooling..."
	@npm install --global prettier
	@echo "Additional tooling..."
	@pre-commit install

# Format Go project
format:
	@go fmt ./...
	@tagalign -fix ./...
	@goimports -w .
	@go clean -i ./...
	@find . -name "*.md" ! -path "./documentation/themes/*" -exec prettier --write {} +

# Upgrade project dependencies
upgrade:
	@echo "Upgrading project dependencies..."
	@go mod tidy
	@go get -u -t ./...

# Run and clean test coverage report
test-coverage:
	@rm -f coverage.out coverage.html
	@echo "Checking test coverage and exporting report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@open coverage.html
	@echo "Cleaning test coverage reports..."

# This uses KO to build a container image.
# - Make sure to run Docker.
# - If there are problems locating the socket, run `docker context ls` 
#   and set DOCKER_HOST (example) DOCKER_HOST=unix:///Users/XYZ/.docker/run/docker.sock 
container-build-local:
	KO_DOCKER_REPO=ko.local ko build .

########################################
### OPEN TOFU

SSH_KEY_PATH ?= $(HOME)/.ssh/id_equinix_carbonaut_ed25519.pub
PRIVATE_KEY_PATH ?= $(HOME)/.ssh/id_equinix_carbonaut_ed25519
CARBONAUT_NUM_PROJECTS ?= 1
CARBONAUT_VM_COUNT_PROJECTS ?= 1

tf:
	./hack/tofu.bash $(cmd)

tf-init:
	@echo "Initializing OpenTofu configuration..."
	@tofu init

tf-plan:
	$(MAKE) tf cmd=plan

tf-apply:
	$(MAKE) tf cmd=apply

tf-destroy:
	$(MAKE) tf cmd=destroy

tf-connect:
	$(MAKE) tf cmd=connect

# uses ansible playbooks
tf-configure:
	$(MAKE) tf cmd=configure

# runs the stress test for all configured machines
tf-stress-test:
	$(MAKE) tf cmd=stress-test

########################################
### GENERAL

help:
	@echo "Available commands:"
	@echo "  build                  - Build the go binary locally"
	@echo "  verify                 - Run verifications on the project (lint, vet, tests)"
	@echo "  install                - Install project dependencies"
	@echo "  format                 - Format Go files"
	@echo "  upgrade                - Upgrade project dependencies"
	@echo "  test-coverage          - Generate and clean test coverage report"
	@echo "  tf-init                - Initialize OpenTofu configuration"
	@echo "  tf-plan                - Plan OpenTofu changes"
	@echo "  tf-apply               - Apply OpenTofu changes"
	@echo "  tf-destroy             - Destroy the created OpenTofu infrastructure"
	@echo "  tf-connect             - Connect to the created server"
	@echo "  tf-configure           - Setup the machine with required packages using Ansible"
	@echo "  tf-stress-test         - Run stress test script on all configured machines"
	@echo "  container-build-local  - Builds a local container image using KO"
	@echo "  SSH_KEY_PATH           - Current SSH key path: $(SSH_KEY_PATH)"
	@echo "  PRIVATE_KEY_PATH       - Current private SSH key path: $(PRIVATE_KEY_PATH)"
