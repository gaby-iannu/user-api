# Project Constitution

## Project Overview

**Name**: User API  
**Language**: Go 1.22+  
**Architecture**: Clean Architecture  

## Core Principles

### Article 1: Code Quality
- Follow Go idioms and best practices
- Use `gofmt` and `golangci-lint` for code formatting and linting
- Keep functions small and focused (single responsibility)
- Prefer composition over inheritance

### Article 2: Error Handling
- Always handle errors explicitly, never ignore them
- Use `errors.Is` and `errors.As` for error type checking
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Define domain-specific error types

### Article 3: Logging
- Use structured logging with `log/slog`
- Log levels: DEBUG, INFO, WARN, ERROR
- Include correlation IDs for request tracing
- Never log sensitive data (passwords, tokens)

### Article 4: Testing
- Write unit tests for all business logic
- Use table-driven tests
- Aim for >80% code coverage on domain logic
- Use interfaces for testability (dependency injection)

### Article 5: API Design
- Follow REST conventions
- Use consistent error response format
- Validate all inputs at the handler level
- Return appropriate HTTP status codes

### Article 6: Security
- Validate and sanitize all user inputs
- Use parameterized queries (no SQL injection)
- Never expose internal errors to clients
- Follow principle of least privilege

### Article 7: Performance
- Response time < 100ms for single resource operations
- Use connection pooling for database
- Implement pagination for list endpoints
- Avoid N+1 query problems

### Article 8: Dependencies
- Prefer standard library when possible
- Minimize external dependencies
- Pin dependency versions
- Document why each dependency is needed

### Article 9: Documentation
- Document public APIs
- Keep README updated
- Use meaningful commit messages
- Specifications are the source of truth
