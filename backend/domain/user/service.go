package user

import (
	"errors"
	"time"
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, username, password string) (*User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	exists, _ := s.repo.GetByUsername(ctx, username)
	if exists != nil {
		return nil, errors.New("username already exists")
	}

	user := &User{
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}