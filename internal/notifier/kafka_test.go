package notifier

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/giannuccilli/user-api/internal/domain"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

type mockFailedEventRepository struct {
	events []domain.FailedEvent
	saveFn func(ctx context.Context, event *domain.FailedEvent) error
}

func (m *mockFailedEventRepository) Save(ctx context.Context, event *domain.FailedEvent) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, event)
	}
	event.ID = uuid.New()
	m.events = append(m.events, *event)
	return nil
}

func (m *mockFailedEventRepository) List(ctx context.Context, limit, offset int) ([]domain.FailedEvent, int, error) {
	return m.events, len(m.events), nil
}

func (m *mockFailedEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestKafkaNotifier_SavesToDLQ_AfterRetries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	mockRepo := &mockFailedEventRepository{}

	n := &KafkaNotifier{
		writer:          nil,
		logger:          testLogger(),
		failedEventRepo: mockRepo,
	}

	event := domain.UserEvent{
		EventID:   uuid.New(),
		EventType: domain.EventTypeUserCreated,
		Timestamp: time.Now().UTC(),
		Data: domain.EventData{
			UserID: uuid.New(),
		},
	}

	payload := []byte(`{"test": "payload"}`)
	lastErr := errors.New("kafka connection refused")

	n.saveToDLQ(context.Background(), event, payload, lastErr)

	if len(mockRepo.events) != 1 {
		t.Fatalf("Expected 1 event in DLQ, got %d", len(mockRepo.events))
	}

	saved := mockRepo.events[0]
	if saved.EventID != event.EventID {
		t.Errorf("EventID = %v, want %v", saved.EventID, event.EventID)
	}
	if saved.EventType != event.EventType {
		t.Errorf("EventType = %v, want %v", saved.EventType, event.EventType)
	}
	if saved.UserID != event.Data.UserID {
		t.Errorf("UserID = %v, want %v", saved.UserID, event.Data.UserID)
	}
	if saved.Error != "kafka connection refused" {
		t.Errorf("Error = %v, want 'kafka connection refused'", saved.Error)
	}
	if saved.Attempts != maxRetries {
		t.Errorf("Attempts = %v, want %v", saved.Attempts, maxRetries)
	}
}

func TestKafkaNotifier_DLQSaveError_LogsButDoesNotPanic(t *testing.T) {
	mockRepo := &mockFailedEventRepository{
		saveFn: func(ctx context.Context, event *domain.FailedEvent) error {
			return errors.New("database error")
		},
	}

	n := &KafkaNotifier{
		writer:          nil,
		logger:          testLogger(),
		failedEventRepo: mockRepo,
	}

	event := domain.UserEvent{
		EventID:   uuid.New(),
		EventType: domain.EventTypeUserCreated,
		Timestamp: time.Now().UTC(),
		Data: domain.EventData{
			UserID: uuid.New(),
		},
	}

	n.saveToDLQ(context.Background(), event, []byte(`{}`), errors.New("kafka error"))
}
