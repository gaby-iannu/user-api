package notifier

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type NoopNotifier struct {
	logger *slog.Logger
}

func NewNoopNotifier(logger *slog.Logger) *NoopNotifier {
	logger.Info("notifications disabled: KAFKA_BROKERS not configured")
	return &NoopNotifier{logger: logger}
}

func (n *NoopNotifier) NotifyCreated(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (n *NoopNotifier) NotifyUpdated(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (n *NoopNotifier) NotifyDeleted(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (n *NoopNotifier) Close() error {
	return nil
}
