#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
COVERAGE=false
VERBOSE=false
SHORT=false
RACE=true
TIMEOUT="10m"
PACKAGE="./..."

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -s|--short)
            SHORT=true
            shift
            ;;
        --no-race)
            RACE=false
            shift
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -p|--package)
            PACKAGE="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  -c, --coverage    Generate coverage report"
            echo "  -v, --verbose     Verbose output"
            echo "  -s, --short       Run only short tests"
            echo "  --no-race         Disable race detector"
            echo "  -t, --timeout     Test timeout (default: 10m)"
            echo "  -p, --package     Specific package to test (default: ./...)"
            echo "  -h, --help        Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🧪 Running tests...${NC}\n"

# Build test command
TEST_CMD="go test"

# Use gotestsum if available
if command -v gotestsum &> /dev/null; then
    TEST_CMD="gotestsum --format testname --"
    echo -e "${GREEN}✓ Using gotestsum for better output${NC}"
fi

# Add flags
TEST_FLAGS=""
if [ "$VERBOSE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -v"
fi
if [ "$SHORT" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -short"
fi
if [ "$RACE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -race"
fi
TEST_FLAGS="$TEST_FLAGS -timeout=$TIMEOUT"

# Add coverage if requested
if [ "$COVERAGE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -coverprofile=coverage.out -covermode=atomic"
fi

# Run tests
echo -e "${YELLOW}📦 Testing package: $PACKAGE${NC}"
echo -e "${YELLOW}⚙️  Flags: $TEST_FLAGS${NC}\n"

if $TEST_CMD $TEST_FLAGS $PACKAGE; then
    echo -e "\n${GREEN}✅ All tests passed!${NC}"
    
    # Generate coverage report if requested
    if [ "$COVERAGE" = true ]; then
        echo -e "\n${YELLOW}📊 Generating coverage report...${NC}"
        go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' > coverage_percentage.txt
        COVERAGE_PCT=$(cat coverage_percentage.txt)
        rm coverage_percentage.txt
        
        echo -e "${BLUE}Total coverage: ${COVERAGE_PCT}%${NC}"
        
        # Generate HTML report
        go tool cover -html=coverage.out -o coverage.html
        echo -e "${GREEN}✓ HTML coverage report generated: coverage.html${NC}"
        
        # Check coverage threshold
        THRESHOLD=70
        if (( $(echo "$COVERAGE_PCT < $THRESHOLD" | bc -l) )); then
            echo -e "${YELLOW}⚠️  Coverage is below ${THRESHOLD}%${NC}"
        else
            echo -e "${GREEN}✓ Coverage meets threshold (${THRESHOLD}%)${NC}"
        fi
    fi
else
    echo -e "\n${RED}❌ Tests failed!${NC}"
    exit 1
fi

# Run benchmarks if requested
if [[ "$PACKAGE" == *"bench"* ]]; then
    echo -e "\n${YELLOW}🏃 Running benchmarks...${NC}"
    go test -bench=. -benchmem $PACKAGE
fi

# Check for test coverage in specific packages
if [ "$COVERAGE" = true ] && [ "$PACKAGE" = "./..." ]; then
    echo -e "\n${YELLOW}📋 Package coverage breakdown:${NC}"
    go tool cover -func=coverage.out | grep -E "^github.com/sachin-duhan/postal-go" | sort -k3 -nr | head -20
fi