#!/bin/bash

set -e

echo "=== Test Coverage Report ==="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Run tests with coverage
echo -e "${YELLOW}[1/3] Running tests with coverage...${NC}"
go test ./... -coverprofile=coverage.out

# Step 2: Generate HTML report
echo -e "${YELLOW}[2/3] Generating HTML report...${NC}"
go tool cover -html=coverage.out -o coverage.html

# Step 3: Show summary
echo -e "${YELLOW}[3/3] Coverage summary:${NC}"
go tool cover -func=coverage.out | tail -1

echo ""
echo -e "${GREEN}=== Report generated: coverage.html ===${NC}"

# Open in browser
open coverage.html
