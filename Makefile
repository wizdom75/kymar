# Kymar - Database Client Pro - Makefile

# Variables
BINARY_NAME=kymar
MAIN_PATH=./cmd/kymar
BUILD_DIR=.
GO=go
GOFLAGS=-v

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
.PHONY: help
help:
	@echo "Kymar - Database Client Pro - Available Commands:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: Run the application (for development)
.PHONY: run
run:
	@echo "🚀 Running Kymar..."
	@$(GO) run $(GOFLAGS) $(MAIN_PATH)

## build: Build the application binary
.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## build-release: Build optimized binary for release
.PHONY: build-release
build-release:
	@echo "🔨 Building release version..."
	@$(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Release build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## install: Install dependencies
.PHONY: install
install:
	@echo "📦 Installing dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "✅ Dependencies installed"

## clean: Remove build artifacts and binaries
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@rm -f $(BUILD_DIR)/*.exe
	@rm -f $(BUILD_DIR)/*.test
	@rm -rf $(BUILD_DIR)/dist
	@echo "✅ Clean complete"

## test: Run tests
.PHONY: test
test:
	@echo "🧪 Running tests..."
	@$(GO) test -v ./...

## test-coverage: Run tests with coverage report
.PHONY: test-coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "✨ Formatting code..."
	@$(GO) fmt ./...
	@echo "✅ Code formatted"

## vet: Run go vet
.PHONY: vet
vet:
	@echo "🔍 Running go vet..."
	@$(GO) vet ./...
	@echo "✅ Vet complete"

## lint: Run linter (requires golangci-lint)
.PHONY: lint
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Run: brew install golangci-lint"; \
	fi

## check: Run fmt, vet, and build
.PHONY: check
check: fmt vet build
	@echo "✅ All checks passed"

## dev: Run with auto-reload (requires air)
.PHONY: dev
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "⚠️  Air not installed. Run: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to regular run..."; \
		$(MAKE) run; \
	fi

## config-dir: Show config directory location
.PHONY: config-dir
config-dir:
	@echo "📁 Config directory: ~/.kymar/"
	@ls -la ~/.kymar/ 2>/dev/null || echo "Config directory not yet created (will be created on first save)"

## config-clean: Remove saved connections
.PHONY: config-clean
config-clean:
	@echo "🗑️  Removing saved connections..."
	@rm -rf ~/.kymar/
	@echo "✅ Saved connections removed"

## deps-update: Update all dependencies
.PHONY: deps-update
deps-update:
	@echo "⬆️  Updating dependencies..."
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "✅ Dependencies updated"

## info: Display project information
.PHONY: info
info:
	@echo "📊 Project Information"
	@echo "======================"
	@echo "Binary Name:     $(BINARY_NAME)"
	@echo "Main Path:       $(MAIN_PATH)"
	@echo "Go Version:      $$($(GO) version)"
	@echo "Build Directory: $(BUILD_DIR)"
	@echo ""
	@echo "📦 Dependencies:"
	@$(GO) list -m all | grep -v "^github.com/pn/dbclient"

## size: Show binary size
.PHONY: size
size:
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "📏 Binary size:"; \
		ls -lh $(BUILD_DIR)/$(BINARY_NAME) | awk '{print "   " $$5 " - " $$9}'; \
	else \
		echo "⚠️  Binary not found. Run 'make build' first."; \
	fi

