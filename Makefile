.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

ide-setup: ## Installs specific requirements for local development
	curl -sSfL \
    	"https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | \
    	sh -s -- -b $$(go env GOPATH)/bin v1.49.0
	go install gotest.tools/gotestsum@v1.8.2
	pre-commit install

lint: ## Run lint
	golangci-lint run ./...

test: ## Run unit tests
	go test ./...

testwatch: ## Run unit tests in watch mode, re-running tests on each file change
	-gotestsum --format pkgname -- -short ./...
	gotestsum --watch --format pkgname -- -short ./...

build: ## Build the binary
	go build ./cmd/goplicate
