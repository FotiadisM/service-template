# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash
.SHELLFLAGS = -o pipefail -ec

.PHONY: default
default: help

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

## path for all development dependencies
export PATH := $(LOCALBIN):$(PATH)

## Tools Versions
GOLANGCILINT_VERSION ?= 1.52.2
BUF_VERSION ?= 1.15.0

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: go-fmt
go-fmt: ## Run gofumpt against code.
	@gofumpt -w .

.PHONY: go-vet
go-vet: ## Run go vet against code.
	@go vet ./...

.PHONY: go-lint
go-lint: golangci-lint ## Run golangci-lint against code.
	@golangci-lint run

.PHONY: test
test: go-vet go-lint ## Run tests using gotestsum.
	gotestsum --format pkgname-and-test-fails -- -race -count=1 ./...

.PHONY: tools
tools: export GOBIN=$(LOCALBIN)
tools: golangci-lint buf | $(LOCALBIN)  ## Install development tools
	@cat tools/tools.go | grep '_' | awk -F '"' '{print $$2}' | xargs -t go install

##@ Build

.PHONY: build
build: ## Compile the app
	@go build -o $(LOCALBIN)/app ./

##@ Schema

.PHONY: schema
schema: schema-apply schema-generate ## Format, apply database schema and generate models

.PHONY: schema-apply
schema-apply: ## Apply schema changes to local database
	atlas schema apply --env local

.PHONY: schema-gen
schema-gen: ## Generate database schema structs and queries
	sqlc generate

##@ Protobuf

.PHONY: proto
proto: proto-fmt proto-lint proto-gen ## Formats, lints and generates protobuf files.

.PHONY: proto-fmt
proto-fmt: buf ## Format protobuf files with buf.
	@buf format -w

.PHONY: proto-lint
proto-lint: buf ## Lint protobuf files with buf.
	@buf lint ./proto

.PHONY: proto-gen
proto-gen: buf proto-fmt proto-lint ## Generate protobuf files with buf.
	@buf generate ./proto

GOLANGCILINT ?= $(LOCALBIN)/golangci-lint
BUF ?= $(LOCALBIN)/buf

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT)
$(GOLANGCILINT): | $(LOCALBIN)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) v$(GOLANGCILINT_VERSION)

.PHONY: buf
buf: $(BUF)
$(BUF): | $(LOCALBIN)
	@curl -sSfL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" -o "$@" && chmod +x "$@"

##@ Clean

.PHONY: clean
clean: ## Remove compiled binaries and build tools
	-rm -rf $(LOCALBIN)
