# Data Model 002: User Notifications

## Entidades

### UserEvent

Representa un evento de dominio relacionado con un usuario.

```go
type EventType string

const (
    EventTypeUserCreated EventType = "user.created"
    EventTypeUserUpdated EventType = "user.updated"
    EventTypeUserDeleted EventType = "user.deleted"
)

type UserEvent struct {
    EventID   uuid.UUID `json:"eventId"`
    EventType EventType `json:"eventType"`
    Timestamp time.Time `json:"timestamp"`
    Data      EventData `json:"data"`
}

type EventData struct {
    UserID uuid.UUID `json:"userId"`
}
```

## Mensaje Kafka

### Estructura JSON

```json
{
  "eventId": "550e8400-e29b-41d4-a716-446655440000",
  "eventType": "user.created",
  "timestamp": "2026-01-12T19:00:00Z",
  "data": {
    "userId": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

### Campos

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `eventId` | UUID | Identificador único del evento |
| `eventType` | string | Tipo de evento: `user.created`, `user.updated`, `user.deleted` |
| `timestamp` | ISO 8601 | Momento en que ocurrió el evento |
| `data.userId` | UUID | ID del usuario afectado |

## Kafka Message

### Headers

| Header | Valor | Descripción |
|--------|-------|-------------|
| `content-type` | `application/json` | Tipo de contenido |
| `event-type` | `user.created` | Tipo de evento (para filtrado) |

### Key

```
userId (string)
```

El `userId` se usa como key para garantizar que todos los eventos de un mismo usuario vayan a la misma partición, manteniendo el orden.

### Value

```json
{
  "eventId": "...",
  "eventType": "...",
  "timestamp": "...",
  "data": {
    "userId": "..."
  }
}
```

## Interfaz del Notifier

```go
// UserNotifier define la interfaz para notificar eventos de usuario
type UserNotifier interface {
    // NotifyCreated publica un evento de usuario creado
    NotifyCreated(ctx context.Context, userID uuid.UUID) error
    
    // NotifyUpdated publica un evento de usuario actualizado
    NotifyUpdated(ctx context.Context, userID uuid.UUID) error
    
    // NotifyDeleted publica un evento de usuario eliminado
    NotifyDeleted(ctx context.Context, userID uuid.UUID) error
    
    // Close cierra la conexión con el broker
    Close() error
}
```

---

## Dead Letter Queue (DLQ)

### Entidad FailedEvent

Representa un evento que no pudo ser publicado a Kafka después de todos los reintentos.

```go
type FailedEvent struct {
    ID         uuid.UUID `json:"id"`
    EventID    uuid.UUID `json:"eventId"`
    EventType  EventType `json:"eventType"`
    UserID     uuid.UUID `json:"userId"`
    Payload    string    `json:"payload"`    // JSON del evento original
    Error      string    `json:"error"`      // Mensaje de error
    Attempts   int       `json:"attempts"`   // Número de intentos realizados
    CreatedAt  time.Time `json:"createdAt"`
    LastError  time.Time `json:"lastError"`  // Timestamp del último error
}
```

### Schema SQL (PostgreSQL)

```sql
CREATE TABLE IF NOT EXISTS failed_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL,
    payload JSONB NOT NULL,
    error TEXT NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_error TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_failed_events_user_id ON failed_events(user_id);
CREATE INDEX idx_failed_events_event_type ON failed_events(event_type);
CREATE INDEX idx_failed_events_created_at ON failed_events(created_at);
```

### Interfaz FailedEventRepository

```go
// FailedEventRepository define la interfaz para persistir eventos fallidos
type FailedEventRepository interface {
    // Save persiste un evento fallido en la DLQ
    Save(ctx context.Context, event *FailedEvent) error
    
    // List obtiene eventos fallidos con paginación
    List(ctx context.Context, limit, offset int) ([]FailedEvent, int, error)
    
    // Delete elimina un evento fallido (después de reprocesarlo)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

---

## Ejemplos de eventos

### user.created

```json
{
  "eventId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "eventType": "user.created",
  "timestamp": "2026-01-12T19:00:00Z",
  "data": {
    "userId": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

### user.updated

```json
{
  "eventId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
  "eventType": "user.updated",
  "timestamp": "2026-01-12T19:05:00Z",
  "data": {
    "userId": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

### user.deleted

```json
{
  "eventId": "c3d4e5f6-a7b8-9012-cdef-123456789012",
  "eventType": "user.deleted",
  "timestamp": "2026-01-12T19:10:00Z",
  "data": {
    "userId": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```
