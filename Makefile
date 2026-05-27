# Podplane <https://podplane.dev>
# Copyright 2026 Nadrama Pty Ltd
# SPDX-License-Identifier: Apache-2.0

IMAGE ?= ghcr.io/podplane/hello
TAG ?= latest
PORT ?= 8080

.DEFAULT_GOAL := help

.PHONY: help setup fmt lint precommit test build run image push

help: ## Show available targets
	@echo "Usage: make <target>"
	@awk 'BEGIN {FS = ":.*?## "} /^##@/ {printf "\n\033[1m%s\033[0m\n", substr($$0, 5)} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Install git hooks
	@mkdir -p .git/hooks
	@printf '%s\n' '#!/usr/bin/env bash' 'set -eo pipefail' 'echo "Running pre-commit checks..."' 'make precommit' 'echo "Pre-commit checks passed."' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Installed .git/hooks/pre-commit"

fmt: ## Format Go source files
	@go fmt ./...

lint: ## Run linters
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is required but not installed"; exit 1; }
	@golangci-lint run

precommit: ## Check formatting and run linters (read-only)
	@echo "Checking formatting..."
	@UNFORMATTED=$$(gofmt -l . 2>&1); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "The following files need formatting (run 'make fmt'):"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi
	@$(MAKE) lint

test: ## Run tests with race detector
	go test -v -race ./...

build: ## Build the Go binary
	go build ./...

run: ## Run the app locally
	PORT=$(PORT) go run .

image: ## Build the container image
	docker build -f Containerfile -t $(IMAGE):$(TAG) .

push: image ## Push the container image
	docker push $(IMAGE):$(TAG)
