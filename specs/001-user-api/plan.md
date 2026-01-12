# Implementation Plan: User API

**Feature**: 001-user-api  
**Branch**: 001-user-api  
**Date**: 2026-01-12  
**Spec**: [spec.md](./spec.md)

## Technology Stack

| Component      | Technology           | Rationale                                    |
|----------------|----------------------|----------------------------------------------|
| Language       | Go 1.22+             | Performance, concurrency, standard library   |
| Router         | net/http (ServeMux)  | Go 1.22+ routing patterns, no dependencies   |
| Database       | PostgreSQL           | ACID compliance, UUID support, JSON support  |
| DB Driver      | pgx/v5               | Native PostgreSQL driver, connection pooling |
| Logging        | log/slog             | Standard library, structured logging         |
| Validation     | Custom               | Minimal dependencies, full control           |
| UUID           | google/uuid          | Standard UUID generation                     |
| Config         | Environment vars     | 12-factor app compliance                     |

## Architecture

```
cmd/
└── api/
    └── main.go              # Entry point, server setup

internal/
├── domain/
│   ├── user.go              # User entity and repository interface
│   └── errors.go            # Domain errors
├── handler/
│   ├── user.go              # HTTP handlers
│   ├── response.go          # Response helpers
│   └── middleware.go        # Logging, recovery middleware
├── repository/
│   └── postgres/
│       └── user.go          # PostgreSQL implementation
├── service/
│   └── user.go              # Business logic
└── config/
    └── config.go            # Configuration loading

migrations/
└── 001_create_users.sql     # Database migration
```

## Design Decisions

### DD-001: Clean Architecture
- **Decision**: Separar en capas domain, service, repository, handler
- **Rationale**: Testabilidad, mantenibilidad, independencia de frameworks
- **Trade-offs**: Más archivos, pero mejor organización

### DD-002: Repository Pattern
- **Decision**: Usar interface para repository en domain layer
- **Rationale**: Permite cambiar implementación de DB sin afectar lógica
- **Trade-offs**: Indirección adicional

### DD-003: No ORM
- **Decision**: Usar SQL directo con pgx
- **Rationale**: Control total, mejor performance, queries explícitas
- **Trade-offs**: Más código manual para mapping

### DD-004: Validation en Handler
- **Decision**: Validar inputs en la capa de handler
- **Rationale**: Fail fast, errores claros al cliente
- **Trade-offs**: Duplicación potencial con validación de DB

### DD-005: UUID como ID
- **Decision**: Usar UUID v4 generado por la base de datos
- **Rationale**: No expone secuencia, distribuible, estándar
- **Trade-offs**: Mayor tamaño que integer

## API Contracts

### POST /api/v1/users

**Request:**
```http
POST /api/v1/users HTTP/1.1
Content-Type: application/json

{
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Response 201:**
```http
HTTP/1.1 201 Created
Content-Type: application/json
Location: /api/v1/users/{id}

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "status": "active",
  "createdAt": "2026-01-12T15:04:05Z",
  "updatedAt": "2026-01-12T15:04:05Z"
}
```

### GET /api/v1/users

**Request:**
```http
GET /api/v1/users?limit=20&offset=0 HTTP/1.1
```

**Response 200:**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "data": [...],
  "pagination": {
    "total": 100,
    "limit": 20,
    "offset": 0
  }
}
```

### GET /api/v1/users/{id}

**Request:**
```http
GET /api/v1/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
```

**Response 200:**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "status": "active",
  "createdAt": "2026-01-12T15:04:05Z",
  "updatedAt": "2026-01-12T15:04:05Z"
}
```

### PUT /api/v1/users/{id}

**Request:**
```http
PUT /api/v1/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Content-Type: application/json

{
  "firstName": "Johnny"
}
```

**Response 200:**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "firstName": "Johnny",
  "lastName": "Doe",
  "status": "active",
  "createdAt": "2026-01-12T15:04:05Z",
  "updatedAt": "2026-01-12T15:04:05Z"
}
```

### DELETE /api/v1/users/{id}

**Request:**
```http
DELETE /api/v1/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
```

**Response 204:**
```http
HTTP/1.1 204 No Content
```

## Local Development (Docker)

### Start PostgreSQL

```bash
docker compose up -d
```

### Stop PostgreSQL

```bash
docker compose down
```

### Reset Database (delete all data)

```bash
docker compose down -v
docker compose up -d
```

### View Logs

```bash
docker compose logs -f postgres
```

### Connect to Database

```bash
docker exec -it user-api-db psql -U userapi -d userapi
```

## Configuration

| Variable          | Required | Default        | Description              |
|-------------------|----------|----------------|--------------------------|
| PORT              | No       | 8080           | Server port              |
| DATABASE_URL      | Yes      | -              | PostgreSQL connection    |
| LOG_LEVEL         | No       | info           | Logging level            |
| READ_TIMEOUT      | No       | 5s             | HTTP read timeout        |
| WRITE_TIMEOUT     | No       | 10s            | HTTP write timeout       |

### Local DATABASE_URL

```
DATABASE_URL=postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable
```

## Dependencies

```go
require (
    github.com/google/uuid v1.6.0
    github.com/jackc/pgx/v5 v5.5.0
)
```

## Validation Rules

| Field     | Rules                                          |
|-----------|------------------------------------------------|
| email     | Required, valid email format, max 255 chars    |
| firstName | Required, 1-100 chars, trimmed                 |
| lastName  | Required, 1-100 chars, trimmed                 |
| status    | Optional, enum: active, inactive, suspended    |
| id (path) | Required, valid UUID format                    |
| limit     | Optional, 1-100, default 20                    |
| offset    | Optional, >= 0, default 0                      |
