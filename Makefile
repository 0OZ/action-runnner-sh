# Makefile for GitHub Runner Manager

# Binary name
BINARY_NAME=runner-manager

# Build directory
BUILD_DIR=build

# Go parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_MOD=$(GO_CMD) mod

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/runner-manager

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO_MOD) download
	$(GO_MOD) verify

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO_TEST) -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GO_CLEAN)
	rm -rf $(BUILD_DIR)

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Format Go code
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	$(GO_CMD) fmt ./...

# Lint Go code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting Go code..."
	golangci-lint run

# Tidy go modules
.PHONY: tidy
tidy:
	@echo "Tidying go modules..."
	$(GO_MOD) tidy

# Install to system (requires sudo)
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  deps     - Install dependencies"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  run      - Build and run the application"
	@echo "  fmt      - Format Go code"
	@echo "  lint     - Lint Go code (requires golangci-lint)"
	@echo "  tidy     - Tidy go modules"
	@echo "  install  - Install binary to /usr/local/bin (requires sudo)"
	@echo "  help     - Show this help message"

# Default target
.DEFAULT_GOAL := help
