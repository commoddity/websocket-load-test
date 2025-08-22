# WebSocket Load Test Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Build parameters
BINARY_NAME=websocket-load-test
BINARY_PATH=./$(BINARY_NAME)
MAIN_FILE=main.go

# Version and build info
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build clean test test-verbose test-race test-cover lint fmt vet mod-tidy mod-verify help install uninstall

# Default target
all: clean lint test build

# Build the binary
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_FILE)
	@echo "$(GREEN)Build complete: $(BINARY_PATH)$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
	rm -f $(BINARY_NAME)-*
	@echo "$(GREEN)Clean complete$(NC)"

# Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GOTEST) -v ./...

# Run tests with verbose output
test-verbose:
	@echo "$(GREEN)Running tests with verbose output...$(NC)"
	$(GOTEST) -v -count=1 ./...

# Run tests with race detection
test-race:
	@echo "$(GREEN)Running tests with race detection...$(NC)"
	$(GOTEST) -race -v ./...

# Run tests with coverage
test-cover:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Run benchmarks
bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem ./...

# Run linter
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	$(GOLINT) run

# Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOCMD) fmt ./...

# Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GOCMD) vet ./...

# Tidy go modules
mod-tidy:
	@echo "$(GREEN)Tidying go modules...$(NC)"
	$(GOMOD) tidy

# Verify go modules
mod-verify:
	@echo "$(GREEN)Verifying go modules...$(NC)"
	$(GOMOD) verify

# Install the binary to GOPATH/bin
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME)...$(NC)"
	$(GOCMD) install $(LDFLAGS)
	@echo "$(GREEN)Installation complete$(NC)"

# Uninstall the binary from GOPATH/bin
uninstall:
	@echo "$(YELLOW)Uninstalling $(BINARY_NAME)...$(NC)"
	rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)Uninstall complete$(NC)"

# Cross-compile for multiple platforms
build-all: clean
	@echo "$(GREEN)Building for multiple platforms...$(NC)"
	
	# Linux AMD64
	@echo "$(YELLOW)Building for Linux AMD64...$(NC)"
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	
	# Linux ARM64
	@echo "$(YELLOW)Building for Linux ARM64...$(NC)"
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	
	# macOS AMD64
	@echo "$(YELLOW)Building for macOS AMD64...$(NC)"
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	
	# macOS ARM64 (Apple Silicon)
	@echo "$(YELLOW)Building for macOS ARM64...$(NC)"
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	
	# Windows AMD64
	@echo "$(YELLOW)Building for Windows AMD64...$(NC)"
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	
	@echo "$(GREEN)Cross-compilation complete$(NC)"

# Quick development cycle
dev: fmt vet test build

# Comprehensive quality check
quality: fmt vet lint test-race test-cover

# Run a quick smoke test
smoke-test: build
	@echo "$(GREEN)Running smoke test...$(NC)"
	./$(BINARY_NAME) --help > /dev/null
	@echo "$(GREEN)Smoke test passed$(NC)"

# Generate test coverage badge (requires go-coverage-badge)
coverage-badge: test-cover
	@if command -v go-coverage-badge >/dev/null 2>&1; then \
		echo "$(GREEN)Generating coverage badge...$(NC)"; \
		go-coverage-badge -fmt=svg -file=coverage.out > coverage-badge.svg; \
		echo "$(GREEN)Coverage badge generated: coverage-badge.svg$(NC)"; \
	else \
		echo "$(YELLOW)go-coverage-badge not installed. Install with: go install github.com/AlexBeauchemin/go-coverage-badge@latest$(NC)"; \
	fi

# Docker build (if Dockerfile exists)
docker-build:
	@if [ -f Dockerfile ]; then \
		echo "$(GREEN)Building Docker image...$(NC)"; \
		docker build -t $(BINARY_NAME):$(VERSION) .; \
		docker build -t $(BINARY_NAME):latest .; \
		echo "$(GREEN)Docker build complete$(NC)"; \
	else \
		echo "$(YELLOW)Dockerfile not found$(NC)"; \
	fi

# Help target
help:
	@echo "$(GREEN)Available targets:$(NC)"
	@echo "  $(YELLOW)build$(NC)        - Build the binary"
	@echo "  $(YELLOW)clean$(NC)        - Clean build artifacts"
	@echo "  $(YELLOW)test$(NC)         - Run tests"
	@echo "  $(YELLOW)test-verbose$(NC) - Run tests with verbose output"
	@echo "  $(YELLOW)test-race$(NC)    - Run tests with race detection"
	@echo "  $(YELLOW)test-cover$(NC)   - Run tests with coverage report"
	@echo "  $(YELLOW)bench$(NC)        - Run benchmarks"
	@echo "  $(YELLOW)lint$(NC)         - Run linter"
	@echo "  $(YELLOW)fmt$(NC)          - Format code"
	@echo "  $(YELLOW)vet$(NC)          - Run go vet"
	@echo "  $(YELLOW)mod-tidy$(NC)     - Tidy go modules"
	@echo "  $(YELLOW)mod-verify$(NC)   - Verify go modules"
	@echo "  $(YELLOW)install$(NC)      - Install binary to GOPATH/bin"
	@echo "  $(YELLOW)uninstall$(NC)    - Uninstall binary from GOPATH/bin"
	@echo "  $(YELLOW)build-all$(NC)    - Cross-compile for multiple platforms"
	@echo "  $(YELLOW)dev$(NC)          - Quick development cycle (fmt, vet, test, build)"
	@echo "  $(YELLOW)quality$(NC)      - Comprehensive quality check"
	@echo "  $(YELLOW)smoke-test$(NC)   - Run a quick smoke test"
	@echo "  $(YELLOW)coverage-badge$(NC) - Generate test coverage badge"
	@echo "  $(YELLOW)docker-build$(NC) - Build Docker image (if Dockerfile exists)"
	@echo "  $(YELLOW)help$(NC)         - Show this help message"
