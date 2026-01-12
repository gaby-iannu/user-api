#!/bin/bash

set -e

echo "=== User API Startup Script ==="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Start PostgreSQL
echo -e "${YELLOW}[1/3] Starting PostgreSQL...${NC}"
docker-compose up -d

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}[1/3] Waiting for PostgreSQL to be healthy...${NC}"
until docker-compose exec -T postgres pg_isready -U userapi -d userapi > /dev/null 2>&1; do
    sleep 1
done
echo -e "${GREEN}[1/3] PostgreSQL is ready!${NC}"

# Step 2: Set environment variables
echo -e "${YELLOW}[2/3] Setting environment variables...${NC}"
export DATABASE_URL="postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable"
export PORT="${PORT:-8080}"
export LOG_LEVEL="${LOG_LEVEL:-info}"
echo -e "${GREEN}[2/3] Environment configured!${NC}"

# Step 3: Run the API
echo -e "${YELLOW}[3/3] Starting User API on port ${PORT}...${NC}"
echo -e "${GREEN}=== API Ready at http://localhost:${PORT} ===${NC}"
echo ""
echo "Endpoints:"
echo "  POST   /api/v1/users      - Create user"
echo "  GET    /api/v1/users      - List users"
echo "  GET    /api/v1/users/{id} - Get user"
echo "  PUT    /api/v1/users/{id} - Update user"
echo "  DELETE /api/v1/users/{id} - Delete user"
echo ""
echo "Press Ctrl+C to stop"
echo ""

go run ./cmd/api
