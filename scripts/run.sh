#!/bin/bash

set -e

echo "=== User API Startup Script ==="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Start PostgreSQL and Kafka
echo -e "${YELLOW}[1/4] Starting PostgreSQL and Kafka...${NC}"
docker-compose up -d

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}[2/4] Waiting for PostgreSQL to be healthy...${NC}"
until docker-compose exec -T postgres pg_isready -U userapi -d userapi > /dev/null 2>&1; do
    sleep 1
done
echo -e "${GREEN}[2/4] PostgreSQL is ready!${NC}"

# Wait for Kafka to be ready
echo -e "${YELLOW}[3/4] Waiting for Kafka to be healthy...${NC}"
until docker-compose exec -T kafka kafka-broker-api-versions --bootstrap-server localhost:9092 > /dev/null 2>&1; do
    sleep 1
done
echo -e "${GREEN}[3/4] Kafka is ready!${NC}"

# Step 3: Set environment variables
echo -e "${YELLOW}[4/4] Setting environment variables...${NC}"
export DATABASE_URL="postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable"
export KAFKA_BROKERS="localhost:9092"
export KAFKA_TOPIC="${KAFKA_TOPIC:-user-events}"
export PORT="${PORT:-8080}"
export LOG_LEVEL="${LOG_LEVEL:-info}"
echo -e "${GREEN}[4/4] Environment configured!${NC}"

# Step 4: Run the API
echo -e "${GREEN}=== API Ready at http://localhost:${PORT} ===${NC}"
echo ""
echo "Endpoints:"
echo "  POST   /api/v1/users      - Create user"
echo "  GET    /api/v1/users      - List users"
echo "  GET    /api/v1/users/{id} - Get user"
echo "  PUT    /api/v1/users/{id} - Update user"
echo "  DELETE /api/v1/users/{id} - Delete user"
echo ""
echo "Kafka:"
echo "  Topic: ${KAFKA_TOPIC}"
echo "  UI:    http://localhost:8090"
echo ""
echo "Press Ctrl+C to stop"
echo ""

go run ./cmd/api
