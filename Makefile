# Makefile for sql-to-drizzle-schema

# Variables
BINARY_NAME=sql-to-drizzle-schema
GO_VERSION=1.24.1
MAIN_PACKAGE=.
BUILD_DIR=bin
COVERAGE_DIR=coverage

# Default target
.DEFAULT_GOAL := help

# Build the binary
.PHONY: build
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✅ Built $(BINARY_NAME) successfully in $(BUILD_DIR)/"

# Build for multiple platforms
.PHONY: build-all
build-all: ## Build for multiple platforms (Linux, macOS, Windows)
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux amd64
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	# Linux arm64
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	# macOS amd64
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	# macOS arm64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	# Windows amd64
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "✅ Built binaries for all platforms in $(BUILD_DIR)/"

# Install the binary to GOPATH/bin
.PHONY: install
install: ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install $(MAIN_PACKAGE)
	@echo "✅ Installed $(BINARY_NAME) successfully"

# Run tests
.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	go test ./...
	@echo "✅ All tests passed"

# Run tests with verbose output
.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	go test -cover ./...
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "✅ Coverage report generated in $(COVERAGE_DIR)/coverage.html"

# Run tests and open coverage report in browser
.PHONY: test-coverage-view
test-coverage-view: test-coverage ## Run tests with coverage and open report in browser
	@echo "Opening coverage report in browser..."
	@which open >/dev/null 2>&1 && open $(COVERAGE_DIR)/coverage.html || \
	which xdg-open >/dev/null 2>&1 && xdg-open $(COVERAGE_DIR)/coverage.html || \
	echo "Please open $(COVERAGE_DIR)/coverage.html in your browser"

# Run benchmarks
.PHONY: bench
bench: ## Run benchmark tests
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Format code
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Run linter (requires golangci-lint)
.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@which golangci-lint >/dev/null 2>&1 || (echo "❌ golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run
	@echo "✅ Linting completed"

# Install linter
.PHONY: lint-install
lint-install: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ golangci-lint installed"

# Run security scan (requires gosec)
.PHONY: security
security: ## Run security scan (requires gosec)
	@echo "Running security scan..."
	@which gosec >/dev/null 2>&1 || (echo "❌ gosec not found. Install it with: make security-install" && exit 1)
	gosec -exclude=G304 ./...
	@echo "✅ Security scan completed"

# Install security scanner
.PHONY: security-install
security-install: ## Install gosec security scanner
	@echo "Installing gosec..."
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "✅ gosec installed"

# Vet code
.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "✅ Vet completed"

# Download dependencies
.PHONY: deps
deps: ## Download and verify dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	@echo "✅ Dependencies downloaded and verified"

# Tidy dependencies
.PHONY: tidy
tidy: ## Tidy and vendor dependencies
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "✅ Dependencies tidied"

# Clean build artifacts
.PHONY: clean
clean: ## Clean build artifacts and coverage files
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@rm -f *.log
	@rm -f coverage.out
	@rm -f *.ts # Remove any generated test files
	@echo "✅ Cleaned build artifacts"

# Run all checks (format, vet, lint, test)
.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "✅ All checks completed successfully"

# Setup development environment
.PHONY: setup
setup: deps lint-install ## Setup development environment
	@echo "Setting up development environment..."
	@echo "📦 Installing additional tools (optional)..."
	@echo "  - Run 'make security-install' to install gosec security scanner"
	@echo "✅ Development environment setup completed"

# Run the tool with example file
.PHONY: example
example: build ## Build and run with example PostgreSQL file
	@echo "Running example conversion..."
	@$(BUILD_DIR)/$(BINARY_NAME) ./example/postgres/create-table.sql -o example-output.ts
	@echo "✅ Example conversion completed. Output: example-output.ts"

# Run the tool with verbose output for debugging
.PHONY: debug-example
debug-example: build ## Run example with verbose output for debugging
	@echo "Running example with debug output..."
	@$(BUILD_DIR)/$(BINARY_NAME) ./example/postgres/create-table.sql -o debug-output.ts --dialect postgresql
	@cat debug-output.ts
	@echo "✅ Debug example completed"

# Generate documentation
.PHONY: docs
docs: ## Generate Go documentation
	@echo "Starting documentation server..."
	@echo "📚 Documentation available at http://localhost:6060/pkg/github.com/konojunya/sql-to-drizzle-schema/"
	@echo "Press Ctrl+C to stop the server"
	godoc -http=:6060

# Show available targets
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"}; /^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# CI pipeline (what CI should run)
.PHONY: ci
ci: deps fmt vet lint test ## Run CI pipeline (deps, fmt, vet, lint, test)
	@echo "✅ CI pipeline completed successfully"

# Release preparation
.PHONY: release-prep
release-prep: clean ci build-all ## Prepare for release (clean, ci, build-all)
	@echo "✅ Release preparation completed"
	@echo "📦 Binaries available in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/

# Quick development cycle
.PHONY: dev
dev: fmt test build ## Quick development cycle (format, test, build)
	@echo "✅ Development cycle completed"

# Show project status
.PHONY: status
status: ## Show project status and environment info
	@echo "📊 Project Status:"
	@echo "  Go version: $(shell go version)"
	@echo "  Project: $(BINARY_NAME)"
	@echo "  Main package: $(MAIN_PACKAGE)"
	@echo "  Build directory: $(BUILD_DIR)"
	@echo ""
	@echo "📁 Directory structure:"
	@find . -type f -name "*.go" | head -10
	@echo ""
	@echo "📦 Dependencies:"
	@go list -m all | head -5
	@echo ""
	@echo "🧪 Test status:"
	@go test -short ./... 2>/dev/null && echo "  ✅ Tests passing" || echo "  ❌ Tests failing"