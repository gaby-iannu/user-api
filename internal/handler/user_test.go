package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/giannuccilli/user-api/internal/domain"
	"github.com/giannuccilli/user-api/internal/service"
)

type mockUserRepository struct {
	users   map[uuid.UUID]*domain.User
	byEmail map[string]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:   make(map[uuid.UUID]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.New()
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if user, ok := m.byEmail[email]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) List(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	users := make([]domain.User, 0)
	for _, u := range m.users {
		users = append(users, *u)
	}
	return users, len(users), nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return domain.ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.users[id]; !ok {
		return domain.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

type mockNotifier struct{}

func (m *mockNotifier) NotifyCreated(ctx context.Context, userID uuid.UUID) error { return nil }
func (m *mockNotifier) NotifyUpdated(ctx context.Context, userID uuid.UUID) error { return nil }
func (m *mockNotifier) NotifyDeleted(ctx context.Context, userID uuid.UUID) error { return nil }
func (m *mockNotifier) Close() error                                              { return nil }

func setupTestHandler() (*UserHandler, *mockUserRepository) {
	repo := newMockUserRepository()
	svc := service.NewUserService(repo, &mockNotifier{})
	handler := NewUserHandler(svc)
	return handler, repo
}

func TestUserHandler_Create(t *testing.T) {
	handler, _ := setupTestHandler()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid request",
			body:       `{"email":"test@example.com","firstName":"John","lastName":"Doe"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid json",
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing email",
			body:       `{"firstName":"John","lastName":"Doe"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Create() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	handler, repo := setupTestHandler()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}
	repo.users[user.ID] = user

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "existing user",
			id:         user.ID.String(),
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing user",
			id:         uuid.New().String(),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid uuid",
			id:         "invalid-uuid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			rec := httptest.NewRecorder()

			handler.GetByID(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("GetByID() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_List(t *testing.T) {
	handler, repo := setupTestHandler()

	for i := 0; i < 3; i++ {
		user := &domain.User{
			ID:        uuid.New(),
			Email:     "test" + string(rune('0'+i)) + "@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Status:    domain.UserStatusActive,
		}
		repo.users[user.ID] = user
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %v, want %v", rec.Code, http.StatusOK)
	}

	var response domain.UserList
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Data) != 3 {
		t.Errorf("List() len = %v, want 3", len(response.Data))
	}
}

func TestUserHandler_Update(t *testing.T) {
	handler, repo := setupTestHandler()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}
	repo.users[user.ID] = user
	repo.byEmail[user.Email] = user

	tests := []struct {
		name       string
		id         string
		body       string
		wantStatus int
	}{
		{
			name:       "valid update",
			id:         user.ID.String(),
			body:       `{"firstName":"Jane"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing user",
			id:         uuid.New().String(),
			body:       `{"firstName":"Jane"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid uuid",
			id:         "invalid-uuid",
			body:       `{"firstName":"Jane"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			id:         user.ID.String(),
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+tt.id, bytes.NewBufferString(tt.body))
			req.SetPathValue("id", tt.id)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Update(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Update() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	handler, repo := setupTestHandler()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}
	repo.users[user.ID] = user

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "existing user",
			id:         user.ID.String(),
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "non-existing user",
			id:         uuid.New().String(),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid uuid",
			id:         "invalid-uuid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			rec := httptest.NewRecorder()

			handler.Delete(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Delete() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}
