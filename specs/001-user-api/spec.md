# Feature: User API

**Feature Number**: 001  
**Branch**: 001-user-api  
**Status**: Draft  
**Date**: 2026-01-12

## Overview

API REST para gestión de usuarios que permite crear, leer, actualizar y eliminar usuarios del sistema.

## User Stories

### US-001: Crear usuario

**Como** administrador del sistema  
**Quiero** crear un nuevo usuario  
**Para** permitir el acceso de nuevas personas al sistema

**Criterios de aceptación:**
- El sistema debe validar que el email sea único
- El sistema debe validar formato de email válido
- El sistema debe validar que firstName tenga entre 1 y 100 caracteres
- El sistema debe validar que lastName tenga entre 1 y 100 caracteres
- El sistema debe retornar el usuario creado con su ID generado (UUID)
- El sistema debe asignar status "active" por defecto
- El sistema debe asignar timestamps createdAt y updatedAt
- El sistema debe retornar error 400 si los datos son inválidos
- El sistema debe retornar error 409 si el email ya existe

### US-002: Obtener usuario por ID

**Como** consumidor de la API  
**Quiero** obtener los datos de un usuario específico  
**Para** visualizar su información

**Criterios de aceptación:**
- El sistema debe aceptar un UUID válido como parámetro
- El sistema debe retornar el usuario completo si existe
- El sistema debe retornar error 400 si el ID no es un UUID válido
- El sistema debe retornar error 404 si el usuario no existe

### US-003: Listar usuarios

**Como** consumidor de la API  
**Quiero** obtener una lista paginada de usuarios  
**Para** visualizar todos los usuarios del sistema

**Criterios de aceptación:**
- El sistema debe soportar paginación con parámetros `limit` y `offset`
- El límite máximo debe ser 100
- El límite por defecto debe ser 20
- El offset por defecto debe ser 0
- El sistema debe retornar el total de usuarios en la respuesta
- El sistema debe retornar una lista vacía si no hay usuarios

### US-004: Actualizar usuario

**Como** administrador del sistema  
**Quiero** actualizar los datos de un usuario  
**Para** mantener la información actualizada

**Criterios de aceptación:**
- El sistema debe permitir actualización parcial (solo campos enviados)
- El sistema debe validar unicidad de email si se modifica
- El sistema debe validar formato de campos si se envían
- El sistema debe actualizar el timestamp updatedAt
- El sistema debe retornar el usuario actualizado
- El sistema debe retornar error 400 si los datos son inválidos
- El sistema debe retornar error 404 si el usuario no existe
- El sistema debe retornar error 409 si el nuevo email ya existe

### US-005: Eliminar usuario

**Como** administrador del sistema  
**Quiero** eliminar un usuario  
**Para** remover accesos del sistema

**Criterios de aceptación:**
- El sistema debe eliminar permanentemente el usuario (hard delete)
- El sistema debe retornar 204 sin contenido en caso de éxito
- El sistema debe retornar error 400 si el ID no es un UUID válido
- El sistema debe retornar error 404 si el usuario no existe

## API Endpoints

| Método | Endpoint           | Descripción        | Request Body        | Response       |
|--------|--------------------|--------------------|---------------------|----------------|
| POST   | /api/v1/users      | Crear usuario      | CreateUserRequest   | User (201)     |
| GET    | /api/v1/users      | Listar usuarios    | -                   | UserList (200) |
| GET    | /api/v1/users/{id} | Obtener usuario    | -                   | User (200)     |
| PUT    | /api/v1/users/{id} | Actualizar usuario | UpdateUserRequest   | User (200)     |
| DELETE | /api/v1/users/{id} | Eliminar usuario   | -                   | - (204)        |

## Error Response Format

Todas las respuestas de error deben seguir este formato:

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable message",
  "details": ["optional", "additional", "details"]
}
```

### Error Codes

| Code              | HTTP Status | Description                    |
|-------------------|-------------|--------------------------------|
| INVALID_REQUEST   | 400         | Request body o parámetros inválidos |
| INVALID_ID        | 400         | ID no es un UUID válido        |
| USER_NOT_FOUND    | 404         | Usuario no existe              |
| EMAIL_EXISTS      | 409         | Email ya registrado            |
| INTERNAL_ERROR    | 500         | Error interno del servidor     |

## Non-Functional Requirements

- **Performance**: Respuesta < 100ms para operaciones individuales
- **Logging**: Usar structured logging con slog
- **Validation**: Validar todos los inputs antes de procesar
- **Database**: Usar PostgreSQL con connection pooling

## Out of Scope (v1)

- Autenticación/Autorización
- Soft delete
- Auditoría de cambios
- Búsqueda avanzada
- Filtros en listado
- Ordenamiento en listado
