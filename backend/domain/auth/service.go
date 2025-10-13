package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"workup_fitness/domain/user"
)

type Service struct {
	service user.Service
}

func NewService(service user.Service) *Service {
	return &Service{service: service}
}

func (s *Service) Register(ctx context.Context, username, password string) (*user.User, error) {
	if username == "" {
		return nil, errors.Join(ErrMissingField, errors.New("username is required"))
	}
	if password == "" {
		return nil, errors.Join(ErrMissingField, errors.New("password is required"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.service.Create(ctx, username, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (*user.User, error) {
	user, err := s.service.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCreds
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCreds
	}

	return user, nil
}
