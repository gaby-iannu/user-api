package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEventType_Values(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  string
	}{
		{EventTypeUserCreated, "user.created"},
		{EventTypeUserUpdated, "user.updated"},
		{EventTypeUserDeleted, "user.deleted"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			if string(tt.eventType) != tt.expected {
				t.Errorf("EventType = %v, want %v", tt.eventType, tt.expected)
			}
		})
	}
}

func TestUserEvent_JSON(t *testing.T) {
	eventID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	timestamp := time.Date(2026, 1, 12, 19, 0, 0, 0, time.UTC)

	event := UserEvent{
		EventID:   eventID,
		EventType: EventTypeUserCreated,
		Timestamp: timestamp,
		Data: EventData{
			UserID: userID,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal UserEvent: %v", err)
	}

	var decoded UserEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal UserEvent: %v", err)
	}

	if decoded.EventID != eventID {
		t.Errorf("EventID = %v, want %v", decoded.EventID, eventID)
	}
	if decoded.EventType != EventTypeUserCreated {
		t.Errorf("EventType = %v, want %v", decoded.EventType, EventTypeUserCreated)
	}
	if decoded.Data.UserID != userID {
		t.Errorf("Data.UserID = %v, want %v", decoded.Data.UserID, userID)
	}
}

func TestUserEvent_JSONFields(t *testing.T) {
	eventID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	timestamp := time.Date(2026, 1, 12, 19, 0, 0, 0, time.UTC)

	event := UserEvent{
		EventID:   eventID,
		EventType: EventTypeUserCreated,
		Timestamp: timestamp,
		Data: EventData{
			UserID: userID,
		},
	}

	data, _ := json.Marshal(event)
	jsonStr := string(data)

	if !contains(jsonStr, `"eventId"`) {
		t.Error("JSON should contain 'eventId' field")
	}
	if !contains(jsonStr, `"eventType"`) {
		t.Error("JSON should contain 'eventType' field")
	}
	if !contains(jsonStr, `"timestamp"`) {
		t.Error("JSON should contain 'timestamp' field")
	}
	if !contains(jsonStr, `"userId"`) {
		t.Error("JSON should contain 'userId' field")
	}
}

func TestFailedEvent_Struct(t *testing.T) {
	id := uuid.New()
	eventID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	failedEvent := FailedEvent{
		ID:        id,
		EventID:   eventID,
		EventType: EventTypeUserCreated,
		UserID:    userID,
		Payload:   `{"test": "payload"}`,
		Error:     "connection refused",
		Attempts:  3,
		CreatedAt: now,
		LastError: now,
	}

	if failedEvent.ID != id {
		t.Errorf("ID = %v, want %v", failedEvent.ID, id)
	}
	if failedEvent.EventID != eventID {
		t.Errorf("EventID = %v, want %v", failedEvent.EventID, eventID)
	}
	if failedEvent.EventType != EventTypeUserCreated {
		t.Errorf("EventType = %v, want %v", failedEvent.EventType, EventTypeUserCreated)
	}
	if failedEvent.UserID != userID {
		t.Errorf("UserID = %v, want %v", failedEvent.UserID, userID)
	}
	if failedEvent.Attempts != 3 {
		t.Errorf("Attempts = %v, want 3", failedEvent.Attempts)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
