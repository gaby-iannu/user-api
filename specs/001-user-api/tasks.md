# Tasks: User API

**Feature**: 001-user-api  
**Date**: 2026-01-12  
**Plan**: [plan.md](./plan.md)

## Task Legend

- `[ ]` Pending
- `[P]` Can be parallelized
- `[X]` Completed

---

## Phase 1: Project Setup

### Task 1.1: Initialize Go module
- [X] Create `go.mod` with module name `github.com/giannuccilli/user-api`
- [X] Set Go version to 1.22+ (1.24.4 installed)

### Task 1.2: Create directory structure [P]
- [X] Create `cmd/api/`
- [X] Create `internal/domain/`
- [X] Create `internal/handler/`
- [X] Create `internal/repository/postgres/`
- [X] Create `internal/service/`
- [X] Create `internal/config/`
- [X] Create `migrations/` (ya existía)

### Task 1.3: Add dependencies
- [X] Add `github.com/google/uuid` v1.6.0
- [X] Add `github.com/jackc/pgx/v5` v5.8.0

---

## Phase 2: Domain Layer

### Task 2.1: Define User entity
- [X] Create `internal/domain/user.go`
- [X] Define `User` struct with all fields
- [X] Define `UserStatus` type with constants
- [X] Define `CreateUserRequest` struct
- [X] Define `UpdateUserRequest` struct
- [X] Define `UserList` struct with pagination

### Task 2.2: Define repository interface
- [X] Define `UserRepository` interface in `internal/domain/user.go`
- [X] Methods: Create, GetByID, GetByEmail, List, Update, Delete

### Task 2.3: Define domain errors
- [X] Create `internal/domain/errors.go`
- [X] Define `ErrUserNotFound`
- [X] Define `ErrEmailExists`
- [X] Define `ErrInvalidInput`

---

## Phase 3: Repository Layer

### Task 3.1: Create database migration
- [X] Create `migrations/001_create_users.sql` (ya existía)
- [X] Define `user_status` enum type
- [X] Define `users` table
- [X] Add indexes

### Task 3.2: Implement PostgreSQL repository
- [X] Create `internal/repository/postgres/user.go`
- [X] Implement `NewUserRepository` constructor
- [X] Implement `Create` method
- [X] Implement `GetByID` method
- [X] Implement `GetByEmail` method
- [X] Implement `List` method with pagination
- [X] Implement `Update` method
- [X] Implement `Delete` method

---

## Phase 4: Service Layer

### Task 4.1: Implement user service
- [X] Create `internal/service/user.go`
- [X] Define `UserService` struct
- [X] Implement `NewUserService` constructor
- [X] Implement `Create` with email uniqueness check
- [X] Implement `GetByID`
- [X] Implement `List`
- [X] Implement `Update` with email uniqueness check
- [X] Implement `Delete`

---

## Phase 5: Handler Layer

### Task 5.1: Create response helpers
- [X] Create `internal/handler/response.go`
- [X] Implement `JSON` response helper
- [X] Implement `Error` response helper
- [X] Define error codes mapping

### Task 5.2: Create middleware [P]
- [X] Create `internal/handler/middleware.go`
- [X] Implement logging middleware
- [X] Implement recovery middleware
- [X] Implement request ID middleware

### Task 5.3: Implement user handlers
- [X] Create `internal/handler/user.go`
- [X] Define `UserHandler` struct
- [X] Implement `NewUserHandler` constructor
- [X] Implement `Create` handler (POST /api/v1/users)
- [X] Implement `GetByID` handler (GET /api/v1/users/{id})
- [X] Implement `List` handler (GET /api/v1/users)
- [X] Implement `Update` handler (PUT /api/v1/users/{id})
- [X] Implement `Delete` handler (DELETE /api/v1/users/{id})

### Task 5.4: Implement input validation
- [X] Validate email format (in service layer)
- [X] Validate firstName length (in service layer)
- [X] Validate lastName length (in service layer)
- [X] Validate status enum (in service layer)
- [X] Validate UUID format (in handler layer)
- [X] Validate pagination params (in service layer)

---

## Phase 6: Configuration & Main

### Task 6.1: Implement configuration
- [X] Create `internal/config/config.go`
- [X] Define `Config` struct
- [X] Implement `Load` function from environment

### Task 6.2: Implement main entry point
- [X] Create `cmd/api/main.go`
- [X] Load configuration
- [X] Initialize database connection
- [X] Initialize repository
- [X] Initialize service
- [X] Initialize handlers
- [X] Setup routes with ServeMux
- [X] Apply middleware
- [X] Start HTTP server
- [X] Implement graceful shutdown

---

## Phase 7: Testing

### Task 7.1: Unit tests - Domain [P]
- [X] Test User entity validation

### Task 7.2: Unit tests - Service [P]
- [X] Test Create user
- [X] Test GetByID
- [X] Test List
- [X] Test Update
- [X] Test Delete
- [X] Test email uniqueness

### Task 7.3: Unit tests - Handler [P]
- [X] Test Create endpoint
- [X] Test GetByID endpoint
- [X] Test List endpoint
- [X] Test Update endpoint
- [X] Test Delete endpoint
- [X] Test validation errors
- [X] Test not found errors

### Task 7.4: Integration tests
- [X] Test full CRUD flow with real database

---

## Phase 8: Documentation

### Task 8.1: Create README
- [X] Project description (con mención a SDD)
- [X] Prerequisites
- [X] Setup instructions
- [X] API documentation
- [X] Environment variables

### Task 8.2: Create Makefile [P]
- [X] `make build`
- [X] `make run`
- [X] `make test`
- [X] `make db-up/db-down/db-reset`
- [X] `make lint`

---

## Execution Order

**Parallel Group 1** (can run together):
- Task 1.1, 1.2, 1.3

**Sequential**:
- Phase 2 (Domain) → Phase 3 (Repository) → Phase 4 (Service) → Phase 5 (Handler) → Phase 6 (Main)

**Parallel Group 2** (after Phase 6):
- Task 7.1, 7.2, 7.3, 8.1, 8.2

**Final**:
- Task 7.4 (Integration tests)
