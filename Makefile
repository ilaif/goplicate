.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

ide-setup: ## Installs specific requirements for local development
	pre-commit install

lint: ## Run lint
	golangci-lint run ./...

test: ## Run unit tests
	go test -short ./...

build: ## Build the binary
	go build ./cmd/goplicate
