package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"workup_fitness/domain/user"
)

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks workup_fitness/domain/auth Service

type UserService interface {
	Create(ctx context.Context, username, passwordHash string) (*user.User, error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
}

type Service interface {
	Register(ctx context.Context, username, password string) (*user.User, error)
	Login(ctx context.Context, username, password string) (*user.User, error)
}

type serviceImpl struct {
	userService UserService
}

func NewService(service UserService) *serviceImpl {
	return &serviceImpl{userService: service}
}

func (s *serviceImpl) Register(ctx context.Context, username, password string) (*user.User, error) {
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

	user, err := s.userService.Create(ctx, username, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *serviceImpl) Login(ctx context.Context, username, password string) (*user.User, error) {
	user, err := s.userService.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCreds
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCreds
	}

	return user, nil
}
