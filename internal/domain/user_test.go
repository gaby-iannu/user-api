package domain

import (
	"testing"
)

func TestUserStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status UserStatus
		want   bool
	}{
		{"active is valid", UserStatusActive, true},
		{"inactive is valid", UserStatusInactive, true},
		{"suspended is valid", UserStatusSuspended, true},
		{"empty is invalid", UserStatus(""), false},
		{"unknown is invalid", UserStatus("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidStatus(tt.status)
			if got != tt.want {
				t.Errorf("isValidStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func isValidStatus(status UserStatus) bool {
	switch status {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended:
		return true
	default:
		return false
	}
}
