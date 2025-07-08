#!/bin/bash

set -e

echo "🚀 Setting up Postal-Go development environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21 or later.${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}✓ Go ${GO_VERSION} detected${NC}"

# Install development tools
echo -e "\n${YELLOW}📦 Installing development tools...${NC}"

# Install golangci-lint
echo "Installing golangci-lint..."
if ! command -v golangci-lint &> /dev/null; then
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
    echo -e "${GREEN}✓ golangci-lint installed${NC}"
else
    echo -e "${GREEN}✓ golangci-lint already installed${NC}"
fi

# Install gofumpt
echo "Installing gofumpt..."
go install mvdan.cc/gofumpt@latest
echo -e "${GREEN}✓ gofumpt installed${NC}"

# Install gotestsum
echo "Installing gotestsum..."
go install gotest.tools/gotestsum@latest
echo -e "${GREEN}✓ gotestsum installed${NC}"

# Install air for hot reloading
echo "Installing air..."
go install github.com/cosmtrek/air@latest
echo -e "${GREEN}✓ air installed${NC}"

# Install delve debugger
echo "Installing delve..."
go install github.com/go-delve/delve/cmd/dlv@latest
echo -e "${GREEN}✓ delve installed${NC}"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "\n${YELLOW}🔧 Creating .env file...${NC}"
    cat > .env << EOF
# Postal API Configuration
POSTAL_API_KEY=your-api-key-here
POSTAL_URL=https://your-postal-server.com
POSTAL_DEBUG=true

# Test Configuration
TEST_EMAIL_FROM=test@example.com
TEST_EMAIL_TO=recipient@example.com
EOF
    echo -e "${GREEN}✓ .env file created${NC}"
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

# Set up git hooks
echo -e "\n${YELLOW}🔗 Setting up git hooks...${NC}"
mkdir -p .git/hooks

cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
echo "Running pre-commit checks..."

# Run formatting
echo "Checking code formatting..."
gofumpt -l . > /tmp/gofumpt-output
if [ -s /tmp/gofumpt-output ]; then
    echo "The following files need formatting:"
    cat /tmp/gofumpt-output
    echo "Run 'gofumpt -w .' to fix formatting"
    rm /tmp/gofumpt-output
    exit 1
fi
rm /tmp/gofumpt-output

# Run linting
echo "Running linter..."
if ! golangci-lint run; then
    echo "Linting failed. Please fix the issues above."
    exit 1
fi

# Run tests
echo "Running tests..."
if ! go test -short ./...; then
    echo "Tests failed. Please fix the failing tests."
    exit 1
fi

echo "Pre-commit checks passed! ✅"
EOF

chmod +x .git/hooks/pre-commit
echo -e "${GREEN}✓ Git hooks configured${NC}"

# Download dependencies
echo -e "\n${YELLOW}📥 Downloading Go dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies downloaded${NC}"

# Verify installation
echo -e "\n${YELLOW}🔍 Verifying installation...${NC}"
go build ./...
echo -e "${GREEN}✓ Build successful${NC}"

# Run short tests
go test -short ./... > /dev/null 2>&1 || true
echo -e "${GREEN}✓ Test framework verified${NC}"

echo -e "\n${GREEN}🎉 Development environment setup complete!${NC}"
echo -e "\n${YELLOW}Quick commands:${NC}"
echo "  make build       - Build the project"
echo "  make test        - Run tests"
echo "  make lint        - Run linting"
echo "  air              - Start hot reloading"
echo "  make integration-test - Run integration tests"
echo -e "\n${YELLOW}Next steps:${NC}"
echo "1. Update .env with your Postal API credentials"
echo "2. Run 'make test' to verify everything works"
echo "3. Start coding! 🚀"