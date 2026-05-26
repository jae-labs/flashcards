.PHONY: help build test lint clean install run dev

# Variables
BINARY_NAME=flashcards
GO=go
GOFLAGS=-v
LDFLAGS=-s -w
BUILD_DIR=build
CMD_DIR=cmd/flashcards

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@cd $(CMD_DIR) && CGO_ENABLED=1 $(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o ../../$(BUILD_DIR)/$(BINARY_NAME)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux AMD64
	@cd $(CMD_DIR) && CC=gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o ../../$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	# macOS AMD64
	@cd $(CMD_DIR) && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o ../../$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	# macOS ARM64
	@cd $(CMD_DIR) && CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o ../../$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	@echo "Build complete for all platforms"

test: ## Run tests
	@echo "Running tests..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage
	@$(GO) tool cover -html=coverage.out

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run --timeout=5m

lint-fix: ## Run linter and fix issues
	@echo "Running linter with auto-fix..."
	@golangci-lint run --fix --timeout=5m

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@echo "Clean complete"

install: build ## Install the binary
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Install complete"

run: build ## Build and run
	@./$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Run in development mode (with live reload if available)
	@cd $(CMD_DIR) && CGO_ENABLED=1 $(GO) run .

generate: ## Generate flashcards from markdown
	@./$(BUILD_DIR)/$(BINARY_NAME) generate --path ./

review: ## Review flashcards
	@./$(BUILD_DIR)/$(BINARY_NAME) review

admin: ## Open admin interface
	@./$(BUILD_DIR)/$(BINARY_NAME) admin

fmt: ## Format code
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@goimports -w -local flashcards .

vet: ## Run go vet
	@echo "Running go vet..."
	@$(GO) vet ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying modules..."
	@$(GO) mod tidy

mod-download: ## Download go modules
	@echo "Downloading modules..."
	@$(GO) mod download

deps: mod-download ## Install dependencies
	@echo "Installing development dependencies..."
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Dependencies installed"
