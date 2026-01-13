# Spec 002: User Notifications

## Overview

Sistema de notificaciones que informa eventos de usuarios a un sistema de mensajería (Kafka) para que otros servicios puedan reaccionar a cambios en el dominio de usuarios.

## User Stories

### US-001: Notificar alta de usuario
**Como** sistema externo  
**Quiero** recibir una notificación cuando se crea un usuario  
**Para** poder sincronizar mi base de datos o ejecutar procesos de onboarding

**Criterios de aceptación:**
- Se publica un evento `user.created` en Kafka
- El evento contiene el `id` del usuario creado
- El evento se publica después de que el usuario se persiste exitosamente
- Si falla la publicación, se registra el error pero no se revierte la creación

### US-002: Notificar modificación de usuario
**Como** sistema externo  
**Quiero** recibir una notificación cuando se modifica un usuario  
**Para** poder actualizar mi copia local de los datos

**Criterios de aceptación:**
- Se publica un evento `user.updated` en Kafka
- El evento contiene el `id` del usuario modificado
- El evento se publica después de que la actualización se persiste exitosamente
- Si falla la publicación, se registra el error pero no se revierte la actualización

### US-003: Notificar borrado de usuario
**Como** sistema externo  
**Quiero** recibir una notificación cuando se elimina un usuario  
**Para** poder limpiar datos relacionados en mi sistema

**Criterios de aceptación:**
- Se publica un evento `user.deleted` en Kafka
- El evento contiene el `id` del usuario eliminado
- El evento se publica después de que el usuario se elimina exitosamente
- Si falla la publicación, se registra el error pero no se revierte la eliminación

---

## Eventos

### Topic

```
user-events
```

### Tipos de eventos

| Evento | Descripción | Trigger |
|--------|-------------|---------|
| `user.created` | Usuario creado | POST /api/v1/users |
| `user.updated` | Usuario modificado | PUT /api/v1/users/{id} |
| `user.deleted` | Usuario eliminado | DELETE /api/v1/users/{id} |

### Estructura del evento

```json
{
  "eventId": "uuid",
  "eventType": "user.created | user.updated | user.deleted",
  "timestamp": "2026-01-12T19:00:00Z",
  "data": {
    "userId": "uuid"
  }
}
```

---

## Requisitos No Funcionales

### RNF-001: Resiliencia
- La falla en la publicación de eventos NO debe afectar la operación principal
- Retry automático con backoff exponencial (3 intentos: 1s, 2s, 4s)
- Si todos los reintentos fallan, persistir evento en tabla `failed_events`
- Los errores de publicación deben ser logueados con nivel ERROR
- Se debe incluir el `userId` y `eventType` en el log de error

### RNF-001.1: Dead Letter Queue (DLQ)
- Eventos fallidos se persisten en tabla `failed_events`
- Cada evento fallido incluye: evento original, error, timestamp, intentos
- Endpoint para reprocesar eventos fallidos (futuro)
- Job periódico para reintentar eventos fallidos (futuro)

### RNF-002: Observabilidad
- Cada evento publicado debe ser logueado con nivel INFO
- Métricas de eventos publicados/fallidos (futuro)

### RNF-003: Configuración
- El broker de Kafka debe ser configurable via variable de entorno
- El topic debe ser configurable via variable de entorno
- La funcionalidad de notificaciones debe poder deshabilitarse

### RNF-004: Ordenamiento
- Los eventos de un mismo usuario deben usar el `userId` como partition key
- Esto garantiza ordenamiento de eventos por usuario

---

## Alcance

### En alcance
- Publicación de eventos a Kafka
- Eventos: created, updated, deleted
- Logging de eventos publicados y errores
- Configuración via variables de entorno
- Retry con backoff exponencial
- Persistencia de eventos fallidos (DLQ en PostgreSQL)

### Fuera de alcance
- Consumo de eventos
- Eventos de otros dominios
- Job automático de reprocesamiento de DLQ
- Endpoint de reprocesamiento manual
