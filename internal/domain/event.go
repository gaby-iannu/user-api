package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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

type FailedEvent struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"eventId"`
	EventType EventType `json:"eventType"`
	UserID    uuid.UUID `json:"userId"`
	Payload   string    `json:"payload"`
	Error     string    `json:"error"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"createdAt"`
	LastError time.Time `json:"lastError"`
}

type UserNotifier interface {
	NotifyCreated(ctx context.Context, userID uuid.UUID) error
	NotifyUpdated(ctx context.Context, userID uuid.UUID) error
	NotifyDeleted(ctx context.Context, userID uuid.UUID) error
	Close() error
}

type FailedEventRepository interface {
	Save(ctx context.Context, event *FailedEvent) error
	List(ctx context.Context, limit, offset int) ([]FailedEvent, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
