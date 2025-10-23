package auth

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Creating auth service...")
	defer log.Info().Msg("Created auth service")
	return &serviceImpl{userService: service}
}

func (s *serviceImpl) Register(ctx context.Context, username, password string) (*user.User, error) {
	log.Info().Msgf("Register user with username %s", username)

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

	log.Info().Msgf("Registered user with username %s", username)

	return user, nil
}

func (s *serviceImpl) Login(ctx context.Context, username, password string) (*user.User, error) {
	log.Info().Msgf("Logging in user with username %s", username)

	user, err := s.userService.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCreds
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password))
	if err != nil {
		return nil, ErrInvalidCreds
	}

	log.Info().Msgf("Logged in user with username %s", username)

	return user, nil
}
