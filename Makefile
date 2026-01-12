.PHONY: build run test test-unit test-integration coverage lint clean db-up db-down db-reset help

# Variables
BINARY_NAME=user-api
DATABASE_URL?=postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## build: Build the application binary
build:
	go build -o bin/$(BINARY_NAME) ./cmd/api

## run: Run the application (requires DB)
run: db-up
	@sleep 2
	DATABASE_URL=$(DATABASE_URL) go run ./cmd/api

## test: Run all tests
test:
	go test ./...

## test-unit: Run unit tests only (no DB required)
test-unit:
	go test ./internal/domain/... ./internal/service/... ./internal/handler/...

## test-integration: Run integration tests (requires DB)
test-integration: db-up
	@sleep 2
	go test ./internal/repository/postgres/... -v

## coverage: Generate test coverage report
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | tail -1
	@echo "Report generated: coverage.html"

## lint: Run linter
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

## clean: Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

## db-up: Start PostgreSQL container
db-up:
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker-compose exec -T postgres pg_isready -U userapi -d userapi > /dev/null 2>&1; do sleep 1; done
	@echo "PostgreSQL is ready!"

## db-down: Stop PostgreSQL container
db-down:
	docker-compose down

## db-reset: Reset database (delete all data)
db-reset:
	docker-compose down -v
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker-compose exec -T postgres pg_isready -U userapi -d userapi > /dev/null 2>&1; do sleep 1; done
	@echo "Database reset complete!"

## db-logs: Show PostgreSQL logs
db-logs:
	docker-compose logs -f postgres

## db-shell: Connect to PostgreSQL shell
db-shell:
	docker exec -it user-api-db psql -U userapi -d userapi
