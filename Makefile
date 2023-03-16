#!/usr/bin/env bash

.PHONY: all clean build test lint run help
.DEFAULT_GOAL := help

DOCKER_IMAGE_NAME="0xerr0r/mailcatcher"
BINARY_NAME=mailcatcher
BIN_OUT_DIR=bin

GOLANG_LINT_VERSION=v1.51.2

export PATH=$(shell go env GOPATH)/bin:$(shell echo $$PATH)

all: build test lint ## Build binary (with tests)

clean: ## cleans output directory
	$(shell rm -rf $(BIN_OUT_DIR)/*)

build:  ## Build binary
	go build -v -o $(BIN_OUT_DIR)/$(BINARY_NAME)

test:  ## run tests
	go run github.com/onsi/ginkgo/v2/ginkgo --coverprofile=coverage.txt --covermode=atomic -cover ./...

lint: build ## run golangcli-lint checks
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANG_LINT_VERSION)
	$(shell go env GOPATH)/bin/golangci-lint run

run: build ## Build and run binary
	./$(BIN_OUT_DIR)/$(BINARY_NAME)

fmt: ## gofmt and goimports all go files
	find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

docker-build:  ## Build docker image
	docker build --network=host --tag ${DOCKER_IMAGE_NAME} .

help:  ## Shows help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'