.DEFAULT_GOAL:=help
SHELL:=/bin/bash

.PHONY: all


.PHONY: help
help: ## Displays this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Building & Running:

.PHONY: run
run: ## run binary directly without docker.

.PHONY: run_image
run_image: ## run the 'latest' docker image.
	@echo -e "\033[92m  ---> Running image ... \033[0m"
		docker run \
			--rm \
			-it \
			-p 8443:8443 \
			"openbanking/conformance-dcr:latest"

.PHONY: build
build: ## build the server binary directly.

##@ Dependencies:

.PHONY: init
init: ## initialise.
	@echo -e "\033[92m  ---> Initialising ... \033[0m"
	go mod download
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u golang.org/x/tools/cmd/goimports

##@ Cleanup:

.PHONY: clean
clean: ## run the clean up
	@echo -e "\033[92m  ---> Cleaning ... \033[0m"
	go clean -i -r -cache -testcache -modcache

##@ Testing:

.PHONY: test
test: ## Run the test suite
	go test -bench -cover -benchmem -coverprofile=coverage.out ./...

.PHONY: test-quick ## Run the quick test suite
test-quick:
	go test -short -failfast

.PHONY: fmt
fmt: ## Run gofmt on all go files
	gofmt -w -s .
	goimports -w .
	go clean -i -r -cache -testcache -modcache

.PHONY: lint
lint: ## Basic linting and vetting of code
	golangci-lint run -E golint

.PHONY: full-lint
full-lint: ## Run a more extensive lint suite
	golangci-lint run -E gosec -E unconvert -E dupl -E goconst -E gocyclo -E maligned -E misspell -E unparam -E prealloc -E gochecknoglobals -E nakedret -E gocritic
