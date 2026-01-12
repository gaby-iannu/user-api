# Data Model: User API

**Feature**: 001-user-api  
**Date**: 2026-01-12

## Entities

### User

| Campo     | Tipo        | Constraints                          | Description                   |
|-----------|-------------|--------------------------------------|-------------------------------|
| id        | UUID        | PRIMARY KEY, NOT NULL                | Identificador único           |
| email     | VARCHAR(255)| UNIQUE, NOT NULL                     | Email del usuario             |
| firstName | VARCHAR(100)| NOT NULL, LENGTH 1-100               | Nombre del usuario            |
| lastName  | VARCHAR(100)| NOT NULL, LENGTH 1-100               | Apellido del usuario          |
| status    | ENUM        | NOT NULL, DEFAULT 'active'           | Estado: active/inactive/suspended |
| createdAt | TIMESTAMP   | NOT NULL, DEFAULT CURRENT_TIMESTAMP  | Fecha de creación             |
| updatedAt | TIMESTAMP   | NOT NULL, DEFAULT CURRENT_TIMESTAMP  | Fecha de última actualización |

## Database Schema (PostgreSQL)

```sql
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    status user_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

## DTOs

### CreateUserRequest

```json
{
  "email": "string (required, valid email format)",
  "firstName": "string (required, 1-100 chars)",
  "lastName": "string (required, 1-100 chars)"
}
```

### UpdateUserRequest

```json
{
  "email": "string (optional, valid email format)",
  "firstName": "string (optional, 1-100 chars)",
  "lastName": "string (optional, 1-100 chars)",
  "status": "string (optional, enum: active|inactive|suspended)"
}
```

### User (Response)

```json
{
  "id": "uuid",
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "status": "string",
  "createdAt": "ISO8601 timestamp",
  "updatedAt": "ISO8601 timestamp"
}
```

### UserList (Response)

```json
{
  "data": [User],
  "pagination": {
    "total": "integer",
    "limit": "integer",
    "offset": "integer"
  }
}
```

### Error (Response)

```json
{
  "code": "string",
  "message": "string",
  "details": ["string"] 
}
```

## Field Mappings

| Go Struct | JSON         | Database Column |
|-----------|--------------|-----------------|
| ID        | id           | id              |
| Email     | email        | email           |
| FirstName | firstName    | first_name      |
| LastName  | lastName     | last_name       |
| Status    | status       | status          |
| CreatedAt | createdAt    | created_at      |
| UpdatedAt | updatedAt    | updated_at      |
