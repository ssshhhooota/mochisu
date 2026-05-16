BINARY_NAME := mochisu
BIN_DIR     := bin

.DEFAULT_GOAL := help

.PHONY: help run dev setup install build lint lint-fix

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	DEBUG=1 go run .

dev: ## Start with hot reload via gow
	DEBUG=1 gow run .

setup: ## Install dependencies and development tools
	go mod download
	go install github.com/mitranim/gow@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

install: ## Install the binary into $GOPATH/bin
	go install .

build: ## Build the binary into bin/
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) .

lint: ## Run golangci-lint
	golangci-lint run ./...

lint-fix: ## Auto-fix with golangci-lint
	golangci-lint run --fix ./...
