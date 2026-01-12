# Plan 002: User Notifications

## Stack Tecnológico

| Componente | Tecnología | Justificación |
|------------|------------|---------------|
| Message Broker | Apache Kafka | Escalable, ordenamiento por partición, estándar de industria |
| Cliente Kafka | segmentio/kafka-go | Cliente Go nativo, simple, bien mantenido |
| Serialización | JSON | Simple, legible, compatible con múltiples consumidores |

## Arquitectura

```
┌──────────────────────────────────────────────────────────────────────────┐
│                              User API                                     │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  ┌──────────┐    ┌──────────────┐    ┌───────────────────────┐           │
│  │ Handler  │───▶│   Service    │───▶│     Repository        │           │
│  └──────────┘    └──────┬───────┘    └───────────────────────┘           │
│                         │                                                 │
│                         │ (after success)                                 │
│                         ▼                                                 │
│                  ┌──────────────┐                                         │
│                  │  Notifier    │                                         │
│                  │  (Publisher) │                                         │
│                  └──────┬───────┘                                         │
│                         │                                                 │
│            ┌────────────┼────────────┐                                    │
│            │            │            │                                    │
│            ▼            │            ▼                                    │
│     ┌────────────┐      │     ┌─────────────┐    ┌───────────────────┐   │
│     │   Retry    │      │     │  On Fail    │───▶│  FailedEventRepo  │   │
│     │  (1s,2s,4s)│──────┘     │  (after 3x) │    │   (PostgreSQL)    │   │
│     └────────────┘            └─────────────┘    └───────────────────┘   │
│                                                                           │
└───────────────────────────────────────────────────────────────────────────┘
                          │
                          ▼
                   ┌──────────────┐
                   │    Kafka     │
                   │ user-events  │
                   └──────────────┘
```

## Diseño

### Patrón: Fire and Forget con Retry + DLQ
- El servicio publica el evento después de la operación exitosa
- Si falla la publicación, se reintenta con backoff exponencial (3 intentos)
- Si todos los reintentos fallan, se persiste en tabla `failed_events`
- No se revierte la operación principal
- Esto prioriza la disponibilidad y garantiza que no se pierdan eventos

### Interfaz del Notifier

```go
type UserNotifier interface {
    NotifyCreated(ctx context.Context, userID uuid.UUID) error
    NotifyUpdated(ctx context.Context, userID uuid.UUID) error
    NotifyDeleted(ctx context.Context, userID uuid.UUID) error
}
```

### Implementaciones

1. **KafkaNotifier**: Publica eventos a Kafka
2. **NoopNotifier**: No hace nada (para cuando Kafka no está configurado)

## Estructura de directorios

```
internal/
├── domain/
│   └── event.go              # Definición de eventos e interfaz FailedEventRepository
├── notifier/
│   ├── notifier.go           # Interfaz UserNotifier
│   ├── kafka.go              # Implementación Kafka con retry
│   └── noop.go               # Implementación vacía
├── repository/
│   └── postgres/
│       └── failed_event.go   # Repositorio para eventos fallidos (DLQ)
└── service/
    └── user.go               # Modificar para usar notifier
```

## Variables de entorno

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `KAFKA_BROKERS` | No | - | Lista de brokers separados por coma |
| `KAFKA_TOPIC` | No | user-events | Topic para eventos de usuario |

Si `KAFKA_BROKERS` está vacío, se usa `NoopNotifier`.

## Docker Compose (desarrollo)

```yaml
kafka:
  image: confluentinc/cp-kafka:7.5.0
  ports:
    - "9092:9092"
  environment:
    KAFKA_NODE_ID: 1
    KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
    KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093
    KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
    KAFKA_PROCESS_ROLES: broker,controller
    KAFKA_CONTROLLER_QUORUM_VOTERS: 1@localhost:9093
    KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
    KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qk
```

## Decisiones de diseño

### D-001: Retry con Backoff Exponencial + DLQ
**Decisión**: Reintentar 3 veces con backoff (1s, 2s, 4s), si falla persistir en DLQ  
**Razón**: Maximiza probabilidad de entrega sin perder eventos  
**Trade-off**: Latencia adicional en caso de falla de Kafka (máx 7s)

### D-002: Partition Key = UserID
**Decisión**: Usar `userId` como partition key  
**Razón**: Garantiza orden de eventos por usuario  
**Trade-off**: Posible desbalance si hay usuarios muy activos

### D-003: Inyección de dependencia
**Decisión**: Inyectar `UserNotifier` en el servicio  
**Razón**: Facilita testing y permite cambiar implementación  
**Trade-off**: Más código de inicialización

### D-004: NoopNotifier por defecto
**Decisión**: Si no hay configuración de Kafka, usar NoopNotifier  
**Razón**: La API debe funcionar sin Kafka en desarrollo  
**Trade-off**: Silencioso, podría ocultar errores de configuración
