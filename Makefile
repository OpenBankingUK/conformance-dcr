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
LATEST_VERSION      := v1.0.4

# Go build flags:
LD_FLAGS := "-X main.version=${LATEST_VERSION} -X main.commitHash=${COMMIT_HASH} -X 'main.buildTime=${BUILD_TIME}'"

.PHONY: all
all: fmt lint_fix test build e2e build_image

.PHONY: help
help: ## Displays this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Building & Running:

.PHONY: run
run: build ## run binary directly without docker.
	./dcr -config-path configs/config.json

.PHONY: build
build: ## build the server binary directly.
	@echo -e "\033[92m  ---> Building ... \033[0m"
	go build -ldflags ${LD_FLAGS} -o dcr bitbucket.org/openbankingteam/conformance-dcr/cmd/cli

.PHONY: build_image
build_image: ## build the docker image. Use available args IMAGE_TAG=v1.x.y, ENABLE_IMAGE_SIGNING=1
	@echo -e "\033[92m  ---> Building image ... \033[0m"
	@# We could enable parallel builds for multi-staged builds with `DOCKER_BUILDKIT=1`
	@# See: https://github.com/moby/moby/pull/37151
	@#DOCKER_BUILDKIT=1
	@export DOCKER_CONTENT_TRUST=${ENABLE_IMAGE_SIGNING}
	docker build ${DOCKER_BUILD_ARGS} -t "openbanking/conformance-dcr:${IMAGE_TAG}" .

##@ Dependencies:

.PHONY: tools
tools: ## install go tools (goimports, golangci-lint)
	@echo -e "\033[92m  ---> Installing Go Tools ... \033[0m"
	go get -u golang.org/x/tools/cmd/goimports
	@printf "%b" "\033[93m" "  ---> Installing golangci-lint@v1.16.0 (https://github.com/golangci/golangci-lint) ... " "\033[0m" "\n"
	curl -sfL "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh" | sh -s -- -b $(shell go env GOPATH)/bin v1.21.0

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
e2e: build ## Run the test suite
	./dcr -config-path configs/config.json > run.out || true
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
	golangci-lint run --config ./.golangci.yml ./...

.PHONY: lint_fix
lint_fix: ## Basic linting and vetting of code with fix option enabled
	@printf "%b" "\033[93m" "  ---> Linting with fix enabled ... " "\033[0m" "\n"
	golangci-lint run --fix --config ./.golangci.yml ./...

.PHONY: pre_commit
pre_commit: fmt lint build test e2e
pre_commit: ## pre-commit checks
	@echo -e "\033[92m  ---> pre-commit ... \033[0m"
	go test -cover -count=1 ./...
