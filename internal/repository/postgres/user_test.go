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

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://userapi:userapi123@localhost:5432/userapi?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	testPool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		os.Exit(1)
	}

	if err := testPool.Ping(ctx); err != nil {
		os.Exit(1)
	}

	code := m.Run()

	testPool.Close()
	os.Exit(code)
}

func cleanupTestData(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to cleanup test data: %v", err)
	}
}

func TestUserRepository_Create(t *testing.T) {
	if testPool == nil {
		t.Skip("Database not available")
	}
	cleanupTestData(t)

	repo := NewUserRepository(testPool)
	ctx := context.Background()

	user := &domain.User{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if user.ID == uuid.Nil {
		t.Error("Create() did not set user ID")
	}
	if user.CreatedAt.IsZero() {
		t.Error("Create() did not set CreatedAt")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("Create() did not set UpdatedAt")
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	if testPool == nil {
		t.Skip("Database not available")
	}
	cleanupTestData(t)

	repo := NewUserRepository(testPool)
	ctx := context.Background()

	user1 := &domain.User{
		Email:     "duplicate@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}

	err := repo.Create(ctx, user1)
	if err != nil {
		t.Fatalf("First Create() error = %v", err)
	}

	user2 := &domain.User{
		Email:     "duplicate@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}

	err = repo.Create(ctx, user2)
	if err != domain.ErrEmailExists {
		t.Errorf("Second Create() error = %v, want %v", err, domain.ErrEmailExists)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	if testPool == nil {
		t.Skip("Database not available")
	}
	cleanupTestData(t)

	repo := NewUserRepository(testPool)
	ctx := context.Background()

	user := &domain.User{
		Email:     "getbyid@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}
	repo.Create(ctx, user)

	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if found.ID != user.ID {
		t.Errorf("GetByID() ID = %v, want %v", found.ID, user.ID)
	}
	if found.Email != user.Email {
		t.Errorf("GetByID() Email = %v, want %v", found.Email, user.Email)
	}

	_, err = repo.GetByID(ctx, uuid.New())
	if err != domain.ErrUserNotFound {
		t.Errorf("GetByID() for non-existing error = %v, want %v", err, domain.ErrUserNotFound)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	if testPool == nil {
		t.Skip("Database not available")
	}
	cleanupTestData(t)

	repo := NewUserRepository(testPool)
	ctx := context.Background()

	user := &domain.User{
		Email:     "getbyemail@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}
	repo.Create(ctx, user)

	found, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("GetByEmail() error = %v", err)
	}
	if found.Email != user.Email {
		t.Errorf("GetByEmail() Email = %v, want %v", found.Email, user.Email)
	}

	_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	if err != domain.ErrUserNotFound {
		t.Errorf("GetByEmail() for non-existing error = %v, want %v", err, domain.ErrUserNotFound)
	}
}

func TestUserRepository_List(t *testing.T) {
	if testPool == nil {
		t.Skip("Database not available")
	}
	cleanupTestData(t)

	repo := NewUserRepository(testPool)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		user := &domain.User{
			Email:     "list" + string(rune('0'+i)) + "@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Status:    domain.UserStatusActive,
		}
		repo.Create(ctx, user)
	}

	users, total, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(users) != 5 {
		t.Errorf("List() len = %v, want 5", len(users))
	}
	if total != 5 {
		t.Errorf("List() total = %v, want 5", total)
	}

	users, total, err = repo.List(ctx, 2, 0)
	if err != nil {
		t.Fatalf("List() with limit error = %v", err)
	}
	if len(users) != 2 {
		t.Errorf("List() with limit len = %v, want 2", len(users))
	}
	if total != 5 {
		t.Errorf("List() with limit total = %v, want 5", total)
	}

	users, _, err = repo.List(ctx, 2, 4)
	if err != nil {
		t.Fatalf("List() with offset error = %v", err)
	}
	if len(users) != 1 {
		t.Errorf("List() with offset len = %v, want 1", len(users))
	}
}

func TestUserRepository_Update(t *testing.T) {
	if testPool == nil {
		t.Skip("Database not available")
	}
	cleanupTestData(t)

	repo := NewUserRepository(testPool)
	ctx := context.Background()

	user := &domain.User{
		Email:     "update@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}
	repo.Create(ctx, user)

	user.FirstName = "Jane"
	user.Status = domain.UserStatusInactive

	err := repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	found, _ := repo.GetByID(ctx, user.ID)
	if found.FirstName != "Jane" {
		t.Errorf("Update() FirstName = %v, want Jane", found.FirstName)
	}
	if found.Status != domain.UserStatusInactive {
		t.Errorf("Update() Status = %v, want %v", found.Status, domain.UserStatusInactive)
	}

	nonExistent := &domain.User{
		ID:        uuid.New(),
		Email:     "nonexistent@example.com",
		FirstName: "Test",
		LastName:  "Test",
		Status:    domain.UserStatusActive,
	}
	err = repo.Update(ctx, nonExistent)
	if err != domain.ErrUserNotFound {
		t.Errorf("Update() for non-existing error = %v, want %v", err, domain.ErrUserNotFound)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	if testPool == nil {
		t.Skip("Database not available")
	}
	cleanupTestData(t)

	repo := NewUserRepository(testPool)
	ctx := context.Background()

	user := &domain.User{
		Email:     "delete@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}
	repo.Create(ctx, user)

	err := repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(ctx, user.ID)
	if err != domain.ErrUserNotFound {
		t.Errorf("GetByID() after delete error = %v, want %v", err, domain.ErrUserNotFound)
	}

	err = repo.Delete(ctx, uuid.New())
	if err != domain.ErrUserNotFound {
		t.Errorf("Delete() for non-existing error = %v, want %v", err, domain.ErrUserNotFound)
	}
}

func TestUserRepository_FullCRUDFlow(t *testing.T) {
	if testPool == nil {
		t.Skip("Database not available")
	}
	cleanupTestData(t)

	repo := NewUserRepository(testPool)
	ctx := context.Background()

	user := &domain.User{
		Email:     "crud@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    domain.UserStatusActive,
	}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	t.Logf("Created user with ID: %s", user.ID)

	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if found.Email != user.Email {
		t.Errorf("GetByID() Email mismatch")
	}

	user.FirstName = "Jane"
	if err := repo.Update(ctx, user); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	found, _ = repo.GetByID(ctx, user.ID)
	if found.FirstName != "Jane" {
		t.Errorf("Update() FirstName not updated")
	}

	users, total, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 {
		t.Errorf("List() total = %v, want 1", total)
	}
	if len(users) != 1 {
		t.Errorf("List() len = %v, want 1", len(users))
	}

	if err := repo.Delete(ctx, user.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(ctx, user.ID)
	if err != domain.ErrUserNotFound {
		t.Errorf("GetByID() after delete should return ErrUserNotFound")
	}

	t.Log("Full CRUD flow completed successfully")
}
