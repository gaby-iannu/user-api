package service

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/giannuccilli/user-api/internal/domain"
)

type mockUserRepository struct {
	users    map[uuid.UUID]*domain.User
	byEmail  map[string]*domain.User
	createFn func(ctx context.Context, user *domain.User) error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:   make(map[uuid.UUID]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
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
	total := len(users)

	if offset >= len(users) {
		return []domain.User{}, total, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], total, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return domain.ErrUserNotFound
	}
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	user, ok := m.users[id]
	if !ok {
		return domain.ErrUserNotFound
	}
	delete(m.users, id)
	delete(m.byEmail, user.Email)
	return nil
}

func TestUserService_Create(t *testing.T) {
	tests := []struct {
		name    string
		req     domain.CreateUserRequest
		wantErr error
	}{
		{
			name: "valid user",
			req: domain.CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: nil,
		},
		{
			name: "empty email",
			req: domain.CreateUserRequest{
				Email:     "",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "invalid email format",
			req: domain.CreateUserRequest{
				Email:     "invalid-email",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "empty firstName",
			req: domain.CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "",
				LastName:  "Doe",
			},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "empty lastName",
			req: domain.CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "",
			},
			wantErr: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			svc := NewUserService(repo)

			user, err := svc.Create(context.Background(), tt.req)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
				return
			}

			if user.Email != tt.req.Email {
				t.Errorf("Create() email = %v, want %v", user.Email, tt.req.Email)
			}
			if user.Status != domain.UserStatusActive {
				t.Errorf("Create() status = %v, want %v", user.Status, domain.UserStatusActive)
			}
		})
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	req := domain.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	_, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("First Create() unexpected error = %v", err)
	}

	_, err = svc.Create(context.Background(), req)
	if err != domain.ErrEmailExists {
		t.Errorf("Second Create() error = %v, want %v", err, domain.ErrEmailExists)
	}
}

func TestUserService_GetByID(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	req := domain.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	created, _ := svc.Create(context.Background(), req)

	user, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Errorf("GetByID() unexpected error = %v", err)
	}
	if user.ID != created.ID {
		t.Errorf("GetByID() ID = %v, want %v", user.ID, created.ID)
	}

	_, err = svc.GetByID(context.Background(), uuid.New())
	if err != domain.ErrUserNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, domain.ErrUserNotFound)
	}
}

func TestUserService_List(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	for i := 0; i < 5; i++ {
		req := domain.CreateUserRequest{
			Email:     "test" + string(rune('0'+i)) + "@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}
		svc.Create(context.Background(), req)
	}

	list, err := svc.List(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("List() unexpected error = %v", err)
	}
	if len(list.Data) != 5 {
		t.Errorf("List() len = %v, want 5", len(list.Data))
	}
	if list.Pagination.Total != 5 {
		t.Errorf("List() total = %v, want 5", list.Pagination.Total)
	}

	list, err = svc.List(context.Background(), 2, 0)
	if err != nil {
		t.Errorf("List() unexpected error = %v", err)
	}
	if len(list.Data) != 2 {
		t.Errorf("List() len = %v, want 2", len(list.Data))
	}
}

func TestUserService_Update(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	req := domain.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	created, _ := svc.Create(context.Background(), req)

	newFirstName := "Jane"
	updateReq := domain.UpdateUserRequest{
		FirstName: &newFirstName,
	}

	updated, err := svc.Update(context.Background(), created.ID, updateReq)
	if err != nil {
		t.Errorf("Update() unexpected error = %v", err)
	}
	if updated.FirstName != newFirstName {
		t.Errorf("Update() firstName = %v, want %v", updated.FirstName, newFirstName)
	}

	_, err = svc.Update(context.Background(), uuid.New(), updateReq)
	if err != domain.ErrUserNotFound {
		t.Errorf("Update() error = %v, want %v", err, domain.ErrUserNotFound)
	}
}

func TestUserService_Delete(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	req := domain.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	created, _ := svc.Create(context.Background(), req)

	err := svc.Delete(context.Background(), created.ID)
	if err != nil {
		t.Errorf("Delete() unexpected error = %v", err)
	}

	_, err = svc.GetByID(context.Background(), created.ID)
	if err != domain.ErrUserNotFound {
		t.Errorf("GetByID() after delete error = %v, want %v", err, domain.ErrUserNotFound)
	}

	err = svc.Delete(context.Background(), uuid.New())
	if err != domain.ErrUserNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domain.ErrUserNotFound)
	}
}
