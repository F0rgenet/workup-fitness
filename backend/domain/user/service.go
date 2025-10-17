package user

import (
	"context"
	"database/sql"
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
	found, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return found, nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	found, err := s.repo.GetByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return found, nil
}

func (s *Service) Update(ctx context.Context, user *User) error {
	_, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return ErrUserNotFound
	}
	found, _ := s.GetByUsername(ctx, user.Username)
	if found != nil {
		return ErrAlreadyExists
	}
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}
	return s.repo.Delete(ctx, id)
}
