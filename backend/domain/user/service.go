package user

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, username, passwordHash string) (*User, error) {
	if username == "" {
		return nil, errors.Join(ErrMissingField, errors.New("username is required"))
	}
	if passwordHash == "" {
		return nil, errors.Join(ErrMissingField, errors.New("passwordHash is required"))
	}
	exists, _ := s.repo.GetByUsername(ctx, username)
	if exists != nil {
		return nil, ErrAlreadyExists
	}
	user := &User{Username: username, PasswordHash: passwordHash}
	createdID, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = createdID
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) Update(ctx context.Context, user *User) error {
	_, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return ErrUserNotFound
	}
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}
	if user.ID != id {
		return ErrInvalidPermissions
	}
	return s.repo.Delete(ctx, id)
}
