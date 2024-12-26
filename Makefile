# Project details
APP_NAME := pipeops
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
DIST_DIR := $(BUILD_DIR)/dist

# Go settings
GO := go
GOFLAGS := -mod=readonly
GO_LDFLAGS := -X main.Version=$(VERSION)

# Default target
.PHONY: all
all: build

# Build the CLI
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	$(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME) .

# Run the application
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	$(BIN_DIR)/$(APP_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf $(BUILD_DIR)

# Package the CLI for distribution (e.g., tar.gz)
.PHONY: package
package: build
	@echo "Packaging $(APP_NAME)..."
	mkdir -p $(DIST_DIR)
	tar -czvf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-$(OS)-$(ARCH).tar.gz -C $(BIN_DIR) $(APP_NAME)

# Publish using Goreleaser (requires goreleaser installed)
.PHONY: release
release:
	@echo "Releasing $(APP_NAME)..."
	goreleaser release

# Install the CLI locally
.PHONY: install
install: build
	@echo "Installing $(APP_NAME)..."
	cp $(BIN_DIR)/$(APP_NAME) /usr/local/bin/

# Cross-compile for multiple platforms
.PHONY: cross-compile
cross-compile:
	@echo "Cross-compiling $(APP_NAME) for all supported platforms..."
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all            Build the CLI (default target)"
	@echo "  build          Build the CLI"
	@echo "  run            Run the CLI"
	@echo "  clean          Clean build artifacts"
	@echo "  package        Package the CLI for distribution"
	@echo "  release        Publish the CLI using Goreleaser"
	@echo "  install        Install the CLI locally"
	@echo "  cross-compile  Cross-compile for multiple platforms"
	@echo "  help           Show this help message"
