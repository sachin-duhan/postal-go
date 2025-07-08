.PHONY: build test lint integration-test e2e-test clean setup coverage

# Default target
all: build

# Setup development environment
setup:
	@bash scripts/setup-dev.sh

# Build the project
build:
	@echo "Building project..."
	@go build -v ./...

# Run tests
test:
	@bash scripts/test.sh

# Run tests with coverage
test-coverage:
	@bash scripts/test.sh --coverage

# Run linting
lint:
	@bash scripts/lint.sh

# Run integration tests
integration-test:
	@echo "Starting integration tests..."
	@docker-compose -f tests/integration/docker-compose.yml up -d
	@go test -v ./tests/integration/...
	@docker-compose -f tests/integration/docker-compose.yml down

# Run e2e tests
e2e-test:
	@echo "Running e2e tests..."
	@go test -v ./tests/e2e/...

# Run all tests
test-all: test integration-test e2e-test

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf tmp/
	@rm -f coverage.out coverage.html
	@rm -f postal-cli
	@go clean -cache

# Show test coverage in browser
coverage:
	@bash scripts/test.sh --coverage
	@go tool cover -html=coverage.out

# Install development tools
tools:
	@echo "Installing development tools..."
	@go install mvdan.cc/gofumpt@latest
	@go install gotest.tools/gotestsum@latest
	@go install github.com/cosmtrek/air@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

# Run development server with hot reload
dev:
	@air -c .air.toml

# Format code
fmt:
	@echo "Formatting code..."
	@gofumpt -w .

# Run quick checks before commit
pre-commit: fmt lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  make setup         - Setup development environment"
	@echo "  make build         - Build the project"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make lint          - Run linting"
	@echo "  make fmt           - Format code"
	@echo "  make dev           - Run with hot reload"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make pre-commit    - Run pre-commit checks"
	@echo "  make help          - Show this help"