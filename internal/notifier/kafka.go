package notifier

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"github.com/giannuccilli/user-api/internal/domain"
)

const (
	maxRetries    = 3
	initialDelay  = 1 * time.Second
	backoffFactor = 2
)

type KafkaNotifier struct {
	writer          *kafka.Writer
	logger          *slog.Logger
	failedEventRepo domain.FailedEventRepository
}

func NewKafkaNotifier(brokers, topic string, logger *slog.Logger, failedEventRepo domain.FailedEventRepository) *KafkaNotifier {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(strings.Split(brokers, ",")...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	logger.Info("kafka notifier initialized",
		slog.String("brokers", brokers),
		slog.String("topic", topic),
	)

	return &KafkaNotifier{
		writer:          writer,
		logger:          logger,
		failedEventRepo: failedEventRepo,
	}
}

func (n *KafkaNotifier) NotifyCreated(ctx context.Context, userID uuid.UUID) error {
	return n.publish(ctx, domain.EventTypeUserCreated, userID)
}

func (n *KafkaNotifier) NotifyUpdated(ctx context.Context, userID uuid.UUID) error {
	return n.publish(ctx, domain.EventTypeUserUpdated, userID)
}

func (n *KafkaNotifier) NotifyDeleted(ctx context.Context, userID uuid.UUID) error {
	return n.publish(ctx, domain.EventTypeUserDeleted, userID)
}

func (n *KafkaNotifier) Close() error {
	return n.writer.Close()
}

func (n *KafkaNotifier) publish(ctx context.Context, eventType domain.EventType, userID uuid.UUID) error {
	event := domain.UserEvent{
		EventID:   uuid.New(),
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data: domain.EventData{
			UserID: userID,
		},
	}

	payload, err := json.Marshal(event)
	if err != nil {
		n.logger.Error("failed to marshal event",
			slog.String("event_type", string(eventType)),
			slog.String("user_id", userID.String()),
			slog.String("error", err.Error()),
		)
		return nil
	}

	msg := kafka.Message{
		Key:   []byte(userID.String()),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "content-type", Value: []byte("application/json")},
			{Key: "event-type", Value: []byte(eventType)},
		},
	}

	var lastErr error
	delay := initialDelay

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := n.writer.WriteMessages(ctx, msg)
		if err == nil {
			n.logger.Info("event published",
				slog.String("event_id", event.EventID.String()),
				slog.String("event_type", string(eventType)),
				slog.String("user_id", userID.String()),
			)
			return nil
		}

		lastErr = err
		n.logger.Error("failed to publish event",
			slog.String("event_type", string(eventType)),
			slog.String("user_id", userID.String()),
			slog.Int("attempt", attempt),
			slog.String("error", err.Error()),
		)

		if attempt < maxRetries {
			time.Sleep(delay)
			delay *= backoffFactor
		}
	}

	n.saveToDLQ(ctx, event, payload, lastErr)
	return nil
}

func (n *KafkaNotifier) saveToDLQ(ctx context.Context, event domain.UserEvent, payload []byte, lastErr error) {
	failedEvent := &domain.FailedEvent{
		EventID:   event.EventID,
		EventType: event.EventType,
		UserID:    event.Data.UserID,
		Payload:   string(payload),
		Error:     lastErr.Error(),
		Attempts:  maxRetries,
		CreatedAt: time.Now().UTC(),
		LastError: time.Now().UTC(),
	}

	if err := n.failedEventRepo.Save(ctx, failedEvent); err != nil {
		n.logger.Error("failed to save event to DLQ",
			slog.String("event_id", event.EventID.String()),
			slog.String("event_type", string(event.EventType)),
			slog.String("user_id", event.Data.UserID.String()),
			slog.String("error", err.Error()),
		)
		return
	}

	n.logger.Warn("event saved to DLQ",
		slog.String("event_id", event.EventID.String()),
		slog.String("event_type", string(event.EventType)),
		slog.String("user_id", event.Data.UserID.String()),
		slog.Int("attempts", maxRetries),
	)
}
