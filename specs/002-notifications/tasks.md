# Tasks 002: User Notifications

## Resumen

| Fase | Descripción | Estimación |
|------|-------------|------------|
| 1 | Setup | 20 min |
| 2 | Domain | 20 min |
| 3 | DLQ Repository | 20 min |
| 4 | Notifier (con retry) | 40 min |
| 5 | Integración | 25 min |
| 6 | Testing | 40 min |
| 7 | Documentación | 15 min |

**Total estimado**: ~3 horas

---

## Phase 1: Setup

### Task 1.1: Agregar dependencia kafka-go
- [X] Ejecutar `go get github.com/segmentio/kafka-go`
- [X] Verificar `go.mod`

### Task 1.2: Actualizar docker-compose
- [X] Agregar servicio Kafka (KRaft mode)
- [X] Configurar puertos y variables de entorno
- [X] Verificar que levanta correctamente

### Task 1.3: Actualizar configuración
- [X] Agregar `KafkaBrokers` a `Config`
- [X] Agregar `KafkaTopic` a `Config`
- [X] Cargar desde variables de entorno

---

## Phase 2: Domain

### Task 2.1: Crear entidades de eventos
- [X] Crear `internal/domain/event.go`
- [X] Definir `EventType` (created, updated, deleted)
- [X] Definir `UserEvent` struct
- [X] Definir `EventData` struct
- [X] Definir `FailedEvent` struct

### Task 2.2: Definir interfaces
- [X] Crear interfaz `UserNotifier` en domain
- [X] Métodos: `NotifyCreated`, `NotifyUpdated`, `NotifyDeleted`, `Close`
- [X] Crear interfaz `FailedEventRepository` en domain
- [X] Métodos: `Save`, `List`, `Delete`

---

## Phase 3: DLQ Repository

### Task 3.1: Crear migración para failed_events
- [X] Crear `migrations/002_create_failed_events.sql`
- [X] Crear tabla `failed_events`
- [X] Crear índices

### Task 3.2: Implementar FailedEventRepository
- [X] Crear `internal/repository/postgres/failed_event.go`
- [X] Implementar `Save`
- [X] Implementar `List`
- [X] Implementar `Delete`

---

## Phase 4: Notifier Layer

### Task 4.1: Implementar NoopNotifier
- [X] Crear `internal/notifier/noop.go`
- [X] Implementar interfaz `UserNotifier`
- [X] Loguear que las notificaciones están deshabilitadas (solo al inicio)

### Task 4.2: Implementar KafkaNotifier con Retry
- [X] Crear `internal/notifier/kafka.go`
- [X] Implementar constructor `NewKafkaNotifier`
- [X] Implementar retry con backoff exponencial (1s, 2s, 4s)
- [X] Implementar `NotifyCreated` con retry
- [X] Implementar `NotifyUpdated` con retry
- [X] Implementar `NotifyDeleted` con retry
- [X] Implementar `Close`
- [X] Usar `userId` como partition key
- [X] Si todos los reintentos fallan, guardar en DLQ via `FailedEventRepository`
- [X] Loguear eventos publicados (INFO)
- [X] Loguear errores de publicación (ERROR)
- [X] Loguear eventos guardados en DLQ (WARN)

### Task 4.3: Factory para crear notifier
- [X] Crear función `NewNotifier(cfg, logger, failedEventRepo)`
- [X] Retornar `KafkaNotifier` si `KafkaBrokers` está configurado
- [X] Retornar `NoopNotifier` si no hay configuración

---

## Phase 5: Integración con Service

### Task 5.1: Modificar UserService
- [X] Agregar `UserNotifier` como dependencia
- [X] Modificar constructor `NewUserService`
- [X] Llamar `NotifyCreated` después de crear usuario
- [X] Llamar `NotifyUpdated` después de actualizar usuario
- [X] Llamar `NotifyDeleted` después de eliminar usuario
- [X] No propagar errores de notificación (el notifier maneja DLQ)

### Task 5.2: Actualizar main.go
- [X] Inicializar `FailedEventRepository`
- [X] Inicializar notifier con factory (pasando failedEventRepo)
- [X] Pasar notifier al servicio
- [X] Cerrar notifier en shutdown

---

## Phase 6: Testing

### Task 6.1: Unit tests - Domain
- [X] Test `EventType` values
- [X] Test `UserEvent` serialización JSON
- [X] Test `FailedEvent` struct

### Task 6.2: Unit tests - FailedEventRepository
- [X] Test `Save` evento fallido
- [X] Test `List` con paginación
- [X] Test `Delete`

### Task 6.3: Unit tests - Notifier
- [X] Test `NoopNotifier` no falla
- [X] Test retry logic con mock de Kafka
- [X] Test que después de 3 fallos guarda en DLQ

### Task 6.4: Unit tests - Service
- [X] Crear mock de `UserNotifier`
- [X] Test que `Create` llama a `NotifyCreated`
- [X] Test que `Update` llama a `NotifyUpdated`
- [X] Test que `Delete` llama a `NotifyDeleted`

### Task 6.5: Integration tests [P]
- [X] Test publicación real a Kafka (requiere Kafka)
- [X] Test retry y fallback a DLQ
- [X] Verificar mensaje en topic

---

## Phase 7: Documentación

### Task 7.1: Actualizar README
- [X] Agregar sección de notificaciones
- [X] Documentar variables de entorno de Kafka
- [X] Documentar estrategia de retry y DLQ
- [ ] Agregar ejemplo de consumo de eventos

### Task 7.2: Actualizar scripts
- [X] Modificar `run.sh` para levantar Kafka (opcional)
- [ ] Agregar script para ver eventos en Kafka
- [ ] Agregar script para ver eventos fallidos en DLQ

---

## Execution Order

**Sequential**:
1. Phase 1 (Setup) → Phase 2 (Domain) → Phase 3 (DLQ Repo) → Phase 4 (Notifier) → Phase 5 (Integración)

**Parallel** (después de Phase 5):
- Phase 6 (Testing)
- Phase 7 (Documentación)

---

## Dependencias entre specs

```
001-user-api ◄─── 002-notifications
     │
     └── Modifica: UserService, main.go, docker-compose.yml, migrations/
```

---

## Retry Strategy

```
Intento 1 → Falla → Espera 1s
Intento 2 → Falla → Espera 2s  
Intento 3 → Falla → Espera 4s
Intento 4 → Falla → Guardar en DLQ (failed_events)
```

**Tiempo máximo de bloqueo**: 7 segundos (1+2+4)
