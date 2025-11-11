.PHONY: build build-all clean test lint dev help

# Binary output directory
BIN_DIR := bin
CMD_DIR := cmd/proxy

# Binary names
BINARY_NAME := postgres-proxy

# Go build flags
LDFLAGS := -ldflags="-s -w"

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Run the proxy in development mode
	go run $(CMD_DIR)/main.go

build: ## Build for current platform
	mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go

build-all: build-windows-amd64 build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64 ## Build for all platforms

build-windows-amd64: ## Build for Windows (amd64)
	mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)/main.go

build-darwin-amd64: ## Build for macOS (Intel)
	mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)/main.go

build-darwin-arm64: ## Build for macOS (Apple Silicon)
	mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)/main.go

build-linux-amd64: ## Build for Linux (amd64)
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)/main.go

build-linux-arm64: ## Build for Linux (arm64)
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)/main.go

test: ## Run tests
	go test ./... -v -race -coverprofile=coverage.out

test-coverage: test ## Run tests and open coverage report
	go tool cover -html=coverage.out

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
	rm -f coverage.out

.DEFAULT_GOAL := help
