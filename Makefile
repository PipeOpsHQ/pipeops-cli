# Project details
APP_NAME := pipeops
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
DIST_DIR := $(BUILD_DIR)/dist
DOCKER_IMAGE := ghcr.io/pipeopshq/pipeops-cli

# Go settings
GO := go
GOFLAGS := -mod=readonly
GO_LDFLAGS := -s -w -X github.com/PipeOpsHQ/pipeops-cli/cmd.Version=$(VERSION) -X github.com/PipeOpsHQ/pipeops-cli/cmd.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) -X github.com/PipeOpsHQ/pipeops-cli/cmd.GitCommit=$(shell git rev-parse HEAD)

# Platform detection
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
RESET := \033[0m

# Default target
.PHONY: all
all: build

# Build the CLI
.PHONY: build
build:
	@echo "$(BLUE)Building $(APP_NAME) $(VERSION) for $(GOOS)/$(GOARCH)...$(RESET)"
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME) .
	@echo "$(GREEN)✓ Build complete: $(BIN_DIR)/$(APP_NAME)$(RESET)"

# Build with race detection
.PHONY: build-race
build-race:
	@echo "$(BLUE)Building $(APP_NAME) with race detection...$(RESET)"
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -race -ldflags "$(GO_LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-race .
	@echo "$(GREEN)✓ Build complete: $(BIN_DIR)/$(APP_NAME)-race$(RESET)"

# Run the application
.PHONY: run
run: build
	@echo "$(BLUE)Running $(APP_NAME)...$(RESET)"
	$(BIN_DIR)/$(APP_NAME)

# Run tests
.PHONY: test
test:
	@echo "$(BLUE)Running tests...$(RESET)"
	$(GO) test -v ./...
	@echo "$(GREEN)✓ Tests passed$(RESET)"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GO) test -v -coverprofile=$(BUILD_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "$(GREEN)✓ Coverage report generated: $(BUILD_DIR)/coverage.html$(RESET)"

# Run benchmark tests
.PHONY: bench
bench:
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	$(GO) test -bench=. ./...

# Run linter
.PHONY: lint
lint:
	@echo "$(BLUE)Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Linting complete$(RESET)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(RESET)"; \
	fi

# Format code
.PHONY: fmt
fmt:
	@echo "$(BLUE)Formatting code...$(RESET)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(RESET)"

# Tidy modules
.PHONY: tidy
tidy:
	@echo "$(BLUE)Tidying modules...$(RESET)"
	$(GO) mod tidy
	@echo "$(GREEN)✓ Modules tidied$(RESET)"

# Generate code
.PHONY: generate
generate:
	@echo "$(BLUE)Generating code...$(RESET)"
	$(GO) generate ./...
	@echo "$(GREEN)✓ Code generated$(RESET)"

# Verify dependencies
.PHONY: verify
verify:
	@echo "$(BLUE)Verifying dependencies...$(RESET)"
	$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies verified$(RESET)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "$(BLUE)Cleaning up build artifacts...$(RESET)"
	rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✓ Clean complete$(RESET)"

# Install the CLI locally
.PHONY: install
install: build
	@echo "$(BLUE)Installing $(APP_NAME) to /usr/local/bin...$(RESET)"
	cp $(BIN_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "$(GREEN)✓ Installed successfully$(RESET)"

# Cross-compile for multiple platforms
.PHONY: cross-compile
cross-compile:
	@echo "$(BLUE)Cross-compiling $(APP_NAME) for all supported platforms...$(RESET)"
	@mkdir -p $(DIST_DIR)
	@echo "  → Linux amd64"
	@GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .
	@echo "  → Linux arm64"
	@GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 .
	@echo "  → Darwin amd64"
	@GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .
	@echo "  → Darwin arm64"
	@GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .
	@echo "  → Windows amd64"
	@GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .
	@echo "  → FreeBSD amd64"
	@GOOS=freebsd GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-freebsd-amd64 .
	@echo "$(GREEN)✓ Cross-compilation complete$(RESET)"

# Package the CLI for distribution
.PHONY: package
package: build
	@echo "$(BLUE)Packaging $(APP_NAME) for $(GOOS)/$(GOARCH)...$(RESET)"
	@mkdir -p $(DIST_DIR)
	tar -czvf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz -C $(BIN_DIR) $(APP_NAME)
	@echo "$(GREEN)✓ Package created: $(DIST_DIR)/$(APP_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz$(RESET)"

# Docker build
.PHONY: docker-build
docker-build:
	@echo "$(BLUE)Building Docker image...$(RESET)"
	docker build -t $(DOCKER_IMAGE):latest -t $(DOCKER_IMAGE):$(VERSION) .
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE):latest$(RESET)"

# Docker run
.PHONY: docker-run
docker-run:
	@echo "$(BLUE)Running Docker container...$(RESET)"
	docker run --rm -it $(DOCKER_IMAGE):latest --help

# Docker push
.PHONY: docker-push
docker-push:
	@echo "$(BLUE)Pushing Docker image...$(RESET)"
	docker push $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):$(VERSION)
	@echo "$(GREEN)✓ Docker image pushed$(RESET)"

# Release using Goreleaser
.PHONY: release
release:
	@echo "$(BLUE)Creating release with Goreleaser...$(RESET)"
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --clean; \
		echo "$(GREEN)✓ Release complete$(RESET)"; \
	else \
		echo "$(RED)✗ Goreleaser not installed. Install from https://goreleaser.com/install/$(RESET)"; \
		exit 1; \
	fi

# Dry run release
.PHONY: release-dry-run
release-dry-run:
	@echo "$(BLUE)Running release dry run...$(RESET)"
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
		echo "$(GREEN)✓ Dry run complete$(RESET)"; \
	else \
		echo "$(RED)✗ Goreleaser not installed. Install from https://goreleaser.com/install/$(RESET)"; \
		exit 1; \
	fi

# Tag a new version
.PHONY: tag
tag:
	@echo "$(BLUE)Current version: $(VERSION)$(RESET)"
	@echo "$(YELLOW)Enter new version (e.g., v1.0.0):$(RESET)"
	@read -p "> " NEW_VERSION; \
	if [ -z "$$NEW_VERSION" ]; then \
		echo "$(RED)✗ Version cannot be empty$(RESET)"; \
		exit 1; \
	fi; \
	git tag -a "$$NEW_VERSION" -m "Release $$NEW_VERSION"; \
	echo "$(GREEN)✓ Tagged $$NEW_VERSION$(RESET)"; \
	echo "$(BLUE)Push tag with: git push origin $$NEW_VERSION$(RESET)"

# Development setup
.PHONY: dev-setup
dev-setup:
	@echo "$(BLUE)Setting up development environment...$(RESET)"
	@echo "  → Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/goreleaser/goreleaser@latest
	@echo "  → Installing dependencies..."
	@$(GO) mod download
	@echo "$(GREEN)✓ Development environment ready$(RESET)"

# Check if everything is ready for release
.PHONY: pre-release
pre-release: lint test verify
	@echo "$(BLUE)Running pre-release checks...$(RESET)"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)✗ Working directory is not clean$(RESET)"; \
		git status --short; \
		exit 1; \
	fi
	@echo "$(GREEN)✓ Ready for release$(RESET)"

# Quick development check
.PHONY: check
check: fmt lint test
	@echo "$(GREEN)✓ All checks passed$(RESET)"

# Show project information
.PHONY: info
info:
	@echo "$(BLUE)Project Information$(RESET)"
	@echo "  Name:        $(APP_NAME)"
	@echo "  Version:     $(VERSION)"
	@echo "  Platform:    $(GOOS)/$(GOARCH)"
	@echo "  Go version:  $(shell $(GO) version)"
	@echo "  Build dir:   $(BUILD_DIR)"
	@echo "  Docker image: $(DOCKER_IMAGE)"

# Show help
.PHONY: help
help:
	@echo "$(BLUE)Available targets:$(RESET)"
	@echo ""
	@echo "$(YELLOW)Building:$(RESET)"
	@echo "  build          Build the CLI"
	@echo "  build-race     Build with race detection"
	@echo "  cross-compile  Cross-compile for all platforms"
	@echo "  install        Install the CLI locally"
	@echo ""
	@echo "$(YELLOW)Testing:$(RESET)"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage"
	@echo "  bench          Run benchmarks"
	@echo "  lint           Run linter"
	@echo "  check          Quick development check (fmt, lint, test)"
	@echo ""
	@echo "$(YELLOW)Development:$(RESET)"
	@echo "  run            Run the CLI"
	@echo "  fmt            Format code"
	@echo "  tidy           Tidy modules"
	@echo "  generate       Generate code"
	@echo "  verify         Verify dependencies"
	@echo "  dev-setup      Setup development environment"
	@echo ""
	@echo "$(YELLOW)Docker:$(RESET)"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run Docker container"
	@echo "  docker-push    Push Docker image"
	@echo ""
	@echo "$(YELLOW)Release:$(RESET)"
	@echo "  release        Create release with Goreleaser"
	@echo "  release-dry-run Dry run release"
	@echo "  tag            Tag a new version"
	@echo "  pre-release    Run pre-release checks"
	@echo "  package        Package for distribution"
	@echo ""
	@echo "$(YELLOW)Utility:$(RESET)"
	@echo "  clean          Clean build artifacts"
	@echo "  info           Show project information"
	@echo "  help           Show this help message"
	@echo ""

# Phony targets
.PHONY: all build build-race run test test-coverage bench lint fmt tidy generate verify clean install cross-compile package docker-build docker-run docker-push release release-dry-run tag dev-setup pre-release check info help

# Security-focused build targets
.PHONY: build-secure build-enterprise build-public

# Build with custom configuration injected at build time
build-secure:
	@echo "Building with secure configuration..."
	go build -ldflags "-X 'github.com/PipeOpsHQ/pipeops-cli/internal/config.DefaultClientID=$(CLIENT_ID)' \
		-X 'github.com/PipeOpsHQ/pipeops-cli/internal/config.DefaultAPIURL=$(API_URL)' \
		-X 'github.com/PipeOpsHQ/pipeops-cli/internal/config.DefaultScopes=$(SCOPES)' \
		-X 'github.com/PipeOpsHQ/pipeops-cli/internal/updater.DefaultGitHubRepo=$(GITHUB_REPO)'" \
		-o bin/pipeops .

# Build for enterprise with custom endpoints
build-enterprise:
	@echo "Building enterprise version..."
	go build -ldflags "-X 'github.com/PipeOpsHQ/pipeops-cli/internal/config.DefaultClientID=pipeops_enterprise_client' \
		-X 'github.com/PipeOpsHQ/pipeops-cli/internal/config.DefaultAPIURL=https://enterprise.pipeops.sh' \
		-X 'github.com/PipeOpsHQ/pipeops-cli/internal/updater.DefaultGitHubRepo=PipeOpsHQ/pipeops-cli-enterprise'" \
		-o bin/pipeops-enterprise .

# Build for public release (obfuscated)
build-public:
	@echo "Building for public release..."
	go build -ldflags "-s -w \
		-X 'github.com/PipeOpsHQ/pipeops-cli/internal/config.DefaultClientID=pipeops_public_client' \
		-X 'github.com/PipeOpsHQ/pipeops-cli/internal/config.DefaultAPIURL=https://api.pipeops.sh' \
		-X 'github.com/PipeOpsHQ/pipeops-cli/internal/updater.DefaultGitHubRepo=PipeOpsHQ/pipeops-cli'" \
		-o bin/pipeops .

# Build with stripped symbols for production
build-stripped:
	@echo "Building with stripped symbols..."
	go build -ldflags "-s -w" -trimpath -o bin/pipeops .

# Build with UPX compression (requires UPX to be installed)
build-compressed: build-stripped
	@echo "Compressing binary..."
	@if command -v upx >/dev/null 2>&1; then \
		upx --best --ultra-brute bin/pipeops; \
	else \
		echo "UPX not installed, skipping compression"; \
	fi
