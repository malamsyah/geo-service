SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

# constants
GOLANGCI_VERSION = 1.60.3

all: install tidy vendor ## Initializes all tools

out:
	@mkdir -p out

download: ## Downloads the dependencies
	@go mod download

tidy: ## Cleans up go.mod and go.sum
	@go mod tidy

vendor: ## Vendorizes the dependencies
	@go mod vendor

install: ## Installs all dependencies
	@go mod download

fmt: ## Formats all code with go fmt
	@go fmt ./...

server: fmt build ## Run the app
	./out/bin/server

build: out/bin ## Builds all binaries

GO_BUILD = mkdir -pv "$(@)" && go build -ldflags="-w -s" -o "$(@)" ./...
.PHONY: out/bin
out/bin:
	$(GO_BUILD)

GOLANGCI_LINT = bin/golangci-lint-$(GOLANGCI_VERSION)
$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b bin v$(GOLANGCI_VERSION)
	@mv bin/golangci-lint "$(@)"

lint: fmt $(GOLANGCI_LINT) download ## Lints all code with golangci-lint
	@$(GOLANGCI_LINT) run

lint-reports: out/lint.xml

.PHONY: out/lint.xml
out/lint.xml: $(GOLANGCI_LINT) out download
	@$(GOLANGCI_LINT) run ./... --out-format checkstyle | tee "$(@)"

test: ## Runs all tests
	@go test $(ARGS) ./...

coverage: out/report.json ## Displays coverage per func on cli
	go tool cover -func=out/coverage.out

test-unit: out
	./bin/test_unit
	go tool cover -func=out/coverage.out

test-repository: out
	./bin/test_repository
	go tool cover -func=out/coverage.out

html-coverage: out/report.json ## Displays the coverage results in the browser
	go tool cover -html=out/coverage.out

test-reports: out/report.json

.PHONY: out/report.json
out/report.json: out
	@go test -count 1 ./... -coverprofile=out/coverage.out --json | tee "$(@)"

clean: ## Cleans up everything
	@rm -rf out

ci: lint-reports test-reports ## Executes lint and test and generates reports

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

start-components:
	@cd deployments && docker-compose up -d

stop-components:
	@cd deployments && docker-compose down

.PHONY: generate-mocks
generate-mocks: ## Generates mocks
	./bin/mockgen

.PHONY: docker-build
docker-build: ## Builds the docker image
	docker build -t server .