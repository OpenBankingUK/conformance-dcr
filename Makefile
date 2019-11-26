.DEFAULT_GOAL		:= help
SHELL				:= /bin/bash

# Time the build was started:
# => Fri Aug 30 09:44:14 UTC 2019
BUILD_TIME			:= $(shell date -u)
# Commit hash from git:
# => 227cea43baed3e8be03f8adc8da33bef73cdb377
# => 227cea4
COMMIT_HASH			:= $(shell git rev-list -1 HEAD)
COMMIT_HASH_SHORT	:= $(shell git rev-parse --short HEAD)

# Go build flags:
# => "-X main.version=0.0.1 -X main.commitHash=227cea43baed3e8be03f8adc8da33bef73cdb377 -X 'main.buildTime=Fri Aug 30 09:46:24 UTC 2019'"
LD_FLAGS := "-X main.version=1.0.0 -X main.commitHash=${COMMIT_HASH} -X 'main.buildTime=${BUILD_TIME}'"

.PHONY: all
all: fmt lint test build e2e build_image

.PHONY: help
help: ## Displays this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Building & Running:

.PHONY: run
run: build
run: ## run binary directly without docker.
	./dcr -config-path configs/config.json

.PHONY: run_image
run_image: ## run the 'latest' docker image.
	@echo -e "\033[92m  ---> Running image ... \033[0m"
		docker run \
			--rm \
			-it \
			"openbanking/conformance-dcr:latest"

.PHONY: build
build: ## build the server binary directly.
	@echo -e "\033[92m  ---> Building ... \033[0m"
	go build -ldflags ${LD_FLAGS} -o dcr bitbucket.org/openbankingteam/conformance-dcr/cmd/cli

.PHONY: build_image
build_image: ## build the docker image.
	@echo -e "\033[92m  ---> Building image ... \033[0m"
	docker build -t "openbanking/conformance-dcr:latest" .

##@ Dependencies:

.PHONY: tools
tools: ## install go tools (goimports, golangci-lint)
	@echo -e "\033[92m  ---> Installing Go Tools ... \033[0m"
	go get -u golang.org/x/tools/cmd/goimports
	@printf "%b" "\033[93m" "  ---> Installing golangci-lint@v1.16.0 (https://github.com/golangci/golangci-lint) ... " "\033[0m" "\n"
	curl -sfL "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh" | sh -s -- -b $(shell go env GOPATH)/bin v1.16.0

.PHONY: deps
deps: ## download dependencies
	@echo -e "\033[92m  ---> Downloading dependencies ... \033[0m"
	go mod download

##@ Cleanup:

.PHONY: clean
clean: ## run the clean up
	@echo -e "\033[92m  ---> Cleaning ... \033[0m"
	go clean -i -r -cache -testcache -modcache

##@ Testing:

.PHONY: test
test: ## Run the test suite
	go test -count=1 ./...

.PHONY: e2e
e2e: build test ## Run the test suite
	-./dcr -config-path configs/config.json > run.out
	diff run.out cmd/cli/testdata/ozone.out

.PHONY: code-coverage
code-coverage: ## Generate code coverage
	./coverage.sh

.PHONY: fmt
fmt: ## Run gofmt on all go files
	gofmt -w -s .
	goimports -w .

.PHONY: lint
lint: ## Basic linting and vetting of code
	@printf "%b" "\033[93m" "  ---> Linting ... " "\033[0m" "\n"
	golangci-lint run --fix --config ./.golangci.yml ./...

.PHONY: pre_commit
pre_commit: fmt lint build test
pre_commit: ## pre-commit checks
	@echo -e "\033[92m  ---> pre-commit ... \033[0m"
	go test -cover -count=1 ./...
