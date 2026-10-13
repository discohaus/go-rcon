# Variables
GOCMD := go
GOTEST := $(GOCMD) test
GOCLEAN := $(GOCMD) clean
GOMOD := $(GOCMD) mod
LINTER := golangci-lint

.PHONY: all help gen test test-coverage lint fmt tidy clean cli cli-dev

# Default target
all: lint test

## gen: Generate the embedded CLI icon from the identity submodule
gen:
	rsvg-convert -w 64 -h 64 identity/assets/go-rcon/chip.svg -o internal/cli/go_rcon_chip.png

## help: Shows the available commands
help:
	@echo "Available Commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

## test: Runs all unit tests
test: gen
	$(GOTEST) -v -race ./...

## test-coverage: Runs tests and generates a coverage report in the browser
test-coverage: gen
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

## lint: Runs the linter
lint: gen
	$(LINTER) run

## fmt: Runs go fmt
fmt:
	$(GOCMD) fmt ./...

## tidy: Runs go mod tidy
tidy: gen
	$(GOMOD) tidy

## clean: Runs go clean
clean:
	$(GOCLEAN)
	rm -f coverage.out

cli: gen ## run: Runs the application
	$(GOCMD) run cmd/cli/main.go -i -H localhost -P 25575 -p 1234

cli-dev:
	docker run -it -v ./data:/data -p 25575:25575 -e EULA=TRUE -e RCON_PASSWORD=1234 itzg/minecraft-server
