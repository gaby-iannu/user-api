# User API

API REST para gestión de usuarios construida íntegramente con la metodología **Spec-Driven Development (SDD)**.

## Spec-Driven Development

Este proyecto fue desarrollado siguiendo la metodología [Spec-Driven Development](https://github.blog/ai-and-ml/generative-ai/spec-driven-development-using-markdown-as-a-programming-language-when-building-with-ai/) de GitHub, donde las especificaciones en Markdown son el artefacto principal y el código se genera desde ellas.

### Estructura de especificaciones

```
specs/
└── 001-user-api/
    ├── spec.md          # Especificación funcional (User Stories)
    ├── plan.md          # Plan de implementación técnica
    ├── data-model.md    # Modelo de datos
    └── tasks.md         # Tareas ejecutables
```

## Tech Stack

| Componente | Tecnología |
|------------|------------|
| Lenguaje | Go 1.22+ |
| Router | net/http (ServeMux) |
| Base de datos | PostgreSQL 16 |
| Driver DB | pgx/v5 |
| Message Broker | Apache Kafka |
| Cliente Kafka | segmentio/kafka-go |
| Logging | log/slog |
| Contenedores | Docker Compose |

## Arquitectura

```
cmd/
└── api/
    └── main.go              # Entry point

internal/
├── config/                  # Configuración
├── domain/                  # Entidades y interfaces
├── handler/                 # HTTP handlers
├── notifier/                # Publicación de eventos a Kafka
├── repository/postgres/     # Implementación PostgreSQL
└── service/                 # Lógica de negocio
```

## Requisitos

- Go 1.22+
- Docker y Docker Compose

## Inicio rápido

```bash
# Clonar el repositorio
git clone https://github.com/giannuccilli/user-api.git
cd user-api

# Levantar la aplicación (PostgreSQL + API)
./scripts/run.sh
```

## Scripts disponibles

| Script | Descripción |
|--------|-------------|
| `./scripts/run.sh` | Levanta PostgreSQL, Kafka y ejecuta la API |
| `./scripts/coverage.sh` | Ejecuta tests y genera reporte de cobertura |
| `./scripts/kafka-events.sh` | Ver eventos en Kafka en tiempo real |
| `./scripts/dlq-events.sh` | Ver eventos fallidos en la DLQ |

## Variables de entorno

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `DATABASE_URL` | Sí | - | Connection string de PostgreSQL |
| `PORT` | No | 8080 | Puerto del servidor |
| `LOG_LEVEL` | No | info | Nivel de logging (debug, info, warn, error) |
| `READ_TIMEOUT` | No | 5s | Timeout de lectura HTTP |
| `WRITE_TIMEOUT` | No | 10s | Timeout de escritura HTTP |
| `KAFKA_BROKERS` | No | - | Lista de brokers Kafka (ej: localhost:9092) |
| `KAFKA_TOPIC` | No | user-events | Topic para eventos de usuario |

### Connection string local

```
DATABASE_URL=postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable
```

## API Endpoints

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/api/v1/users` | Crear usuario |
| `GET` | `/api/v1/users` | Listar usuarios (paginado) |
| `GET` | `/api/v1/users/{id}` | Obtener usuario por ID |
| `PUT` | `/api/v1/users/{id}` | Actualizar usuario |
| `DELETE` | `/api/v1/users/{id}` | Eliminar usuario |

## Ejemplos de uso

### Crear usuario

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

### Listar usuarios

```bash
curl http://localhost:8080/api/v1/users?limit=20&offset=0
```

### Obtener usuario

```bash
curl http://localhost:8080/api/v1/users/{id}
```

### Actualizar usuario

```bash
curl -X PUT http://localhost:8080/api/v1/users/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Jane",
    "status": "inactive"
  }'
```

### Eliminar usuario

```bash
curl -X DELETE http://localhost:8080/api/v1/users/{id}
```

## Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests con cobertura
./scripts/coverage.sh

# Solo unit tests (sin DB)
go test ./internal/domain/... ./internal/service/... ./internal/handler/...

# Integration tests (requiere DB)
docker-compose up -d
go test ./internal/repository/postgres/... -v
```

## Notificaciones (Kafka)

La API publica eventos a Kafka cuando se crean, actualizan o eliminan usuarios.

### Eventos

| Evento | Descripción |
|--------|-------------|
| `user.created` | Usuario creado |
| `user.updated` | Usuario actualizado |
| `user.deleted` | Usuario eliminado |

### Estructura del evento

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

### Resiliencia

- **Retry**: 3 intentos con backoff exponencial (1s, 2s, 4s)
- **DLQ**: Si todos los reintentos fallan, el evento se guarda en la tabla `failed_events`
- **No bloquea**: La operación principal (CRUD) nunca falla por errores de Kafka

### Kafka UI

Accede a http://localhost:8090 para ver los mensajes en Kafka.

### Sin Kafka

Si `KAFKA_BROKERS` no está configurado, la API funciona normalmente sin publicar eventos.

### Ejemplo de consumidor (Go)

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/segmentio/kafka-go"
)

type UserEvent struct {
    EventID   string `json:"eventId"`
    EventType string `json:"eventType"`
    Timestamp string `json:"timestamp"`
    Data      struct {
        UserID string `json:"userId"`
    } `json:"data"`
}

func main() {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "user-events",
        GroupID: "my-consumer-group",
    })
    defer reader.Close()

    for {
        msg, err := reader.ReadMessage(context.Background())
        if err != nil {
            log.Fatal(err)
        }

        var event UserEvent
        json.Unmarshal(msg.Value, &event)
        
        fmt.Printf("Event: %s, UserID: %s\n", event.EventType, event.Data.UserID)
    }
}
```

## Desarrollo local

```bash
# Levantar PostgreSQL y Kafka
docker-compose up -d

# Ejecutar la API con Kafka
DATABASE_URL="postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable" \
KAFKA_BROKERS="localhost:9092" \
go run ./cmd/api

# Ejecutar la API sin Kafka
DATABASE_URL="postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable" \
go run ./cmd/api

# Detener servicios
docker-compose down

# Resetear base de datos
docker-compose down -v
docker-compose up -d
```

## Estructura del proyecto

```
user-api/
├── .specify/
│   └── memory/
│       └── constitution.md     # Principios del proyecto
├── api/                        # (futuro) OpenAPI spec
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── domain/
│   ├── handler/
│   ├── repository/
│   │   └── postgres/
│   └── service/
├── migrations/
│   └── 001_create_users.sql
├── scripts/
│   ├── run.sh
│   └── coverage.sh
├── specs/
│   ├── 001-user-api/
│   │   ├── spec.md
│   │   ├── plan.md
│   │   ├── data-model.md
│   │   └── tasks.md
│   └── 002-notifications/
│       ├── spec.md
│       ├── plan.md
│       ├── data-model.md
│       └── tasks.md
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Licencia

MIT
