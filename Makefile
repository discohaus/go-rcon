# Variables
GOCMD := go
GOTEST := $(GOCMD) test
GOCLEAN := $(GOCMD) clean
GOMOD := $(GOCMD) mod
LINTER := golangci-lint

.PHONY: all help test test-coverage lint fmt tidy clean

# Default target
all: lint test

## help: Shows the available commands
help:
	@echo "Available Commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

## test: Runs all unit tests
test:
	$(GOTEST) -v -race ./...

## test-coverage: Runs tests and generates a coverage report in the browser
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

## lint: Runs the linter
lint:
	$(LINTER) run

## fmt: Runs go fmt
fmt:
	$(GOCMD) fmt ./...

## tidy: Runs go mod tidy
tidy:
	$(GOMOD) tidy

## clean: Runs go clean
clean:
	$(GOCLEAN)
	rm -f coverage.out

cli: ## run: Runs the application
	$(GOCMD) run cmd/cli/main.go -H localhost -P 25575 -p 1234

cli-dev:
	docker run -it -v ./data:/data -p 25575:25575 -e EULA=TRUE -e RCON_PASSWORD=1234 itzg/minecraft-server
