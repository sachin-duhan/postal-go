#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Running code quality checks...${NC}\n"

# Check if tools are installed
if ! command -v gofumpt &> /dev/null; then
    echo -e "${RED}❌ gofumpt is not installed. Run './scripts/setup-dev.sh' to install it.${NC}"
    exit 1
fi

if ! command -v golangci-lint &> /dev/null; then
    echo -e "${RED}❌ golangci-lint is not installed. Run './scripts/setup-dev.sh' to install it.${NC}"
    exit 1
fi

# Run gofumpt to check formatting
echo -e "${YELLOW}📝 Checking code formatting with gofumpt...${NC}"
UNFORMATTED=$(gofumpt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo -e "${RED}❌ The following files need formatting:${NC}"
    echo "$UNFORMATTED"
    echo -e "\n${YELLOW}Run 'gofumpt -w .' to fix formatting issues.${NC}"
    FAILED=true
else
    echo -e "${GREEN}✓ Code formatting is correct${NC}"
fi

# Run go mod tidy check
echo -e "\n${YELLOW}📦 Checking go.mod and go.sum...${NC}"
cp go.mod go.mod.backup
cp go.sum go.sum.backup 2>/dev/null || true
go mod tidy
if ! diff -q go.mod go.mod.backup > /dev/null 2>&1 || ! diff -q go.sum go.sum.backup > /dev/null 2>&1; then
    echo -e "${RED}❌ go.mod or go.sum is not tidy. Run 'go mod tidy' to fix.${NC}"
    mv go.mod.backup go.mod
    mv go.sum.backup go.sum 2>/dev/null || true
    FAILED=true
else
    echo -e "${GREEN}✓ go.mod and go.sum are tidy${NC}"
    rm go.mod.backup
    rm go.sum.backup 2>/dev/null || true
fi

# Run golangci-lint
echo -e "\n${YELLOW}🔧 Running golangci-lint...${NC}"
if golangci-lint run --timeout=5m; then
    echo -e "${GREEN}✓ Linting passed${NC}"
else
    echo -e "${RED}❌ Linting failed${NC}"
    FAILED=true
fi

# Check for TODO/FIXME comments
echo -e "\n${YELLOW}📌 Checking for TODO/FIXME comments...${NC}"
TODOS=$(grep -rn "TODO\|FIXME\|HACK\|BUG" --include="*.go" . 2>/dev/null | grep -v vendor || true)
if [ -n "$TODOS" ]; then
    echo -e "${YELLOW}⚠️  Found TODO/FIXME comments:${NC}"
    echo "$TODOS"
fi

# Summary
echo -e "\n${BLUE}📊 Summary:${NC}"
if [ "$FAILED" = true ]; then
    echo -e "${RED}❌ Some checks failed. Please fix the issues above.${NC}"
    exit 1
else
    echo -e "${GREEN}✅ All checks passed!${NC}"
    
    # Offer to auto-fix formatting if running interactively
    if [ -t 0 ] && [ -n "$UNFORMATTED" ]; then
        echo -e "\n${YELLOW}Would you like to automatically fix formatting issues? (y/n)${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            gofumpt -w .
            echo -e "${GREEN}✓ Formatting fixed!${NC}"
        fi
    fi
fi