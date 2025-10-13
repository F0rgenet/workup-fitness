package user

import (
	"errors"
	"time"
	"context"

	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:  username,
		PasswordHash: string(hashedPassword),
		CreatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (*User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*User, error) {
	return s.repo.GetByID(ctx, id)
}