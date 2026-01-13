package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/giannuccilli/user-api/internal/domain"
)

func setupFailedEventTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skipf("Skipping test: could not connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("Skipping test: could not ping database: %v", err)
	}

	return pool
}

func cleanupFailedEvents(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, _ = pool.Exec(ctx, "DELETE FROM failed_events")
}

func TestFailedEventRepository_Save(t *testing.T) {
	pool := setupFailedEventTestDB(t)
	defer pool.Close()
	cleanupFailedEvents(t, pool)

	repo := NewFailedEventRepository(pool)

	event := &domain.FailedEvent{
		EventID:   uuid.New(),
		EventType: domain.EventTypeUserCreated,
		UserID:    uuid.New(),
		Payload:   `{"test": "payload"}`,
		Error:     "connection refused",
		Attempts:  3,
		CreatedAt: time.Now().UTC(),
		LastError: time.Now().UTC(),
	}

	err := repo.Save(context.Background(), event)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if event.ID == uuid.Nil {
		t.Error("Save() should set ID")
	}
}

func TestFailedEventRepository_List(t *testing.T) {
	pool := setupFailedEventTestDB(t)
	defer pool.Close()
	cleanupFailedEvents(t, pool)

	repo := NewFailedEventRepository(pool)

	for i := 0; i < 5; i++ {
		event := &domain.FailedEvent{
			EventID:   uuid.New(),
			EventType: domain.EventTypeUserCreated,
			UserID:    uuid.New(),
			Payload:   `{"test": "payload"}`,
			Error:     "connection refused",
			Attempts:  3,
			CreatedAt: time.Now().UTC(),
			LastError: time.Now().UTC(),
		}
		if err := repo.Save(context.Background(), event); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	events, total, err := repo.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if total != 5 {
		t.Errorf("List() total = %v, want 5", total)
	}

	if len(events) != 5 {
		t.Errorf("List() len = %v, want 5", len(events))
	}

	events, _, err = repo.List(context.Background(), 2, 0)
	if err != nil {
		t.Fatalf("List() with limit error = %v", err)
	}

	if len(events) != 2 {
		t.Errorf("List() with limit len = %v, want 2", len(events))
	}
}

func TestFailedEventRepository_Delete(t *testing.T) {
	pool := setupFailedEventTestDB(t)
	defer pool.Close()
	cleanupFailedEvents(t, pool)

	repo := NewFailedEventRepository(pool)

	event := &domain.FailedEvent{
		EventID:   uuid.New(),
		EventType: domain.EventTypeUserCreated,
		UserID:    uuid.New(),
		Payload:   `{"test": "payload"}`,
		Error:     "connection refused",
		Attempts:  3,
		CreatedAt: time.Now().UTC(),
		LastError: time.Now().UTC(),
	}

	if err := repo.Save(context.Background(), event); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	err := repo.Delete(context.Background(), event.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	events, total, _ := repo.List(context.Background(), 10, 0)
	if total != 0 || len(events) != 0 {
		t.Errorf("Delete() should remove event, got total=%d, len=%d", total, len(events))
	}
}

func TestFailedEventRepository_Delete_NotFound(t *testing.T) {
	pool := setupFailedEventTestDB(t)
	defer pool.Close()
	cleanupFailedEvents(t, pool)

	repo := NewFailedEventRepository(pool)

	err := repo.Delete(context.Background(), uuid.New())
	if err != domain.ErrUserNotFound {
		t.Errorf("Delete() error = %v, want ErrUserNotFound", err)
	}
}
