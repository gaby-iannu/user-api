package notifier

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestNoopNotifier_NotifyCreated(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	n := NewNoopNotifier(logger)

	err := n.NotifyCreated(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("NotifyCreated() error = %v, want nil", err)
	}
}

func TestNoopNotifier_NotifyUpdated(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	n := NewNoopNotifier(logger)

	err := n.NotifyUpdated(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("NotifyUpdated() error = %v, want nil", err)
	}
}

func TestNoopNotifier_NotifyDeleted(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	n := NewNoopNotifier(logger)

	err := n.NotifyDeleted(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("NotifyDeleted() error = %v, want nil", err)
	}
}

func TestNoopNotifier_Close(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	n := NewNoopNotifier(logger)

	err := n.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}
