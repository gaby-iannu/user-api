package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/giannuccilli/user-api/internal/domain"
)

type UserService struct {
	repo     domain.UserRepository
	notifier domain.UserNotifier
}

func NewUserService(repo domain.UserRepository, notifier domain.UserNotifier) *UserService {
	return &UserService{repo: repo, notifier: notifier}
}

func (s *UserService) Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	firstName := strings.TrimSpace(req.FirstName)
	lastName := strings.TrimSpace(req.LastName)

	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validateName(firstName, "firstName"); err != nil {
		return nil, err
	}
	if err := validateName(lastName, "lastName"); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrEmailExists
	}

	user := &domain.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Status:    domain.UserStatusActive,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	s.notifier.NotifyCreated(ctx, user.ID)

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, limit, offset int) (*domain.UserList, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &domain.UserList{
		Data: users,
		Pagination: domain.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*req.Email))
		if err := validateEmail(email); err != nil {
			return nil, err
		}

		if email != user.Email {
			existingUser, err := s.repo.GetByEmail(ctx, email)
			if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
				return nil, err
			}
			if existingUser != nil {
				return nil, domain.ErrEmailExists
			}
			user.Email = email
		}
	}

	if req.FirstName != nil {
		firstName := strings.TrimSpace(*req.FirstName)
		if err := validateName(firstName, "firstName"); err != nil {
			return nil, err
		}
		user.FirstName = firstName
	}

	if req.LastName != nil {
		lastName := strings.TrimSpace(*req.LastName)
		if err := validateName(lastName, "lastName"); err != nil {
			return nil, err
		}
		user.LastName = lastName
	}

	if req.Status != nil {
		if err := validateStatus(*req.Status); err != nil {
			return nil, err
		}
		user.Status = *req.Status
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	s.notifier.NotifyUpdated(ctx, user.ID)

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.notifier.NotifyDeleted(ctx, id)

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return domain.ErrInvalidInput
	}
	if len(email) > 255 {
		return domain.ErrInvalidInput
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return domain.ErrInvalidInput
	}
	return nil
}

func validateName(name, field string) error {
	if name == "" {
		return domain.ErrInvalidInput
	}
	if len(name) > 100 {
		return domain.ErrInvalidInput
	}
	return nil
}

func validateStatus(status domain.UserStatus) error {
	switch status {
	case domain.UserStatusActive, domain.UserStatusInactive, domain.UserStatusSuspended:
		return nil
	default:
		return domain.ErrInvalidInput
	}
}
