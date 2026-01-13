//go:build integration

package notifier

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"github.com/giannuccilli/user-api/internal/domain"
)

const (
	testBrokers = "localhost:9092"
	testTopic   = "user-events-integration"
)

func ensureTopicExists(t *testing.T) {
	t.Helper()

	conn, err := kafka.Dial("tcp", testBrokers)
	if err != nil {
		t.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		t.Fatalf("Failed to get controller: %v", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		t.Fatalf("Failed to connect to controller: %v", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             testTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	_ = controllerConn.CreateTopics(topicConfigs...)
}

func skipIfNoKafka(t *testing.T) {
	t.Helper()
	if os.Getenv("KAFKA_INTEGRATION") != "true" {
		t.Skip("Skipping Kafka integration test. Set KAFKA_INTEGRATION=true to run.")
	}
}

func TestKafkaNotifier_Integration_PublishAndVerify(t *testing.T) {
	skipIfNoKafka(t)
	ensureTopicExists(t)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := &mockFailedEventRepository{}
	notifier := NewKafkaNotifier(testBrokers, testTopic, logger, mockRepo)
	defer notifier.Close()

	ctx := context.Background()

	conn, err := kafka.DialLeader(ctx, "tcp", testBrokers, testTopic, 0)
	if err != nil {
		t.Fatalf("Failed to connect to partition leader: %v", err)
	}
	offset, _ := conn.ReadLastOffset()
	conn.Close()

	userID := uuid.New()
	err = notifier.NotifyCreated(ctx, userID)
	if err != nil {
		t.Fatalf("NotifyCreated() error = %v", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{testBrokers},
		Topic:     testTopic,
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})
	defer reader.Close()

	reader.SetOffset(offset)

	readCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	msg, err := reader.ReadMessage(readCtx)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	if string(msg.Key) != userID.String() {
		t.Errorf("Key = %s, want %s", string(msg.Key), userID.String())
	}

	var event domain.UserEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if event.EventType != domain.EventTypeUserCreated {
		t.Errorf("EventType = %v, want %v", event.EventType, domain.EventTypeUserCreated)
	}
	if event.Data.UserID != userID {
		t.Errorf("UserID = %v, want %v", event.Data.UserID, userID)
	}

	t.Log("Successfully published and verified event in Kafka")
}

func TestKafkaNotifier_Integration_AllEventTypes(t *testing.T) {
	skipIfNoKafka(t)
	ensureTopicExists(t)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := &mockFailedEventRepository{}
	notifier := NewKafkaNotifier(testBrokers, testTopic, logger, mockRepo)
	defer notifier.Close()

	ctx := context.Background()

	tests := []struct {
		name      string
		notify    func(uuid.UUID) error
		eventType domain.EventType
	}{
		{"created", func(id uuid.UUID) error { return notifier.NotifyCreated(ctx, id) }, domain.EventTypeUserCreated},
		{"updated", func(id uuid.UUID) error { return notifier.NotifyUpdated(ctx, id) }, domain.EventTypeUserUpdated},
		{"deleted", func(id uuid.UUID) error { return notifier.NotifyDeleted(ctx, id) }, domain.EventTypeUserDeleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := uuid.New()
			err := tt.notify(userID)
			if err != nil {
				t.Errorf("Notify%s() error = %v", tt.name, err)
			}
		})
	}

	t.Log("All event types published successfully")
}

func TestKafkaNotifier_Integration_DLQOnFailure(t *testing.T) {
	if os.Getenv("KAFKA_INTEGRATION") != "true" {
		t.Skip("Skipping Kafka integration test. Set KAFKA_INTEGRATION=true to run.")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := &mockFailedEventRepository{}

	notifier := NewKafkaNotifier("invalid-broker:9092", testTopic, logger, mockRepo)
	defer notifier.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userID := uuid.New()

	_ = notifier.NotifyCreated(ctx, userID)

	if len(mockRepo.events) != 1 {
		t.Errorf("Expected 1 event in DLQ, got %d", len(mockRepo.events))
	}

	if len(mockRepo.events) > 0 {
		saved := mockRepo.events[0]
		if saved.UserID != userID {
			t.Errorf("DLQ UserID = %v, want %v", saved.UserID, userID)
		}
		if saved.Attempts != maxRetries {
			t.Errorf("DLQ Attempts = %v, want %v", saved.Attempts, maxRetries)
		}
	}
}
