package user

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks workup_fitness/domain/user Service

type Service interface {
	Create(ctx context.Context, username, passwordHash string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) *serviceImpl {
	log.Info().Msg("Creating user service...")
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) Create(ctx context.Context, username, passwordHash string) (*User, error) {
	defer log.Info().Msgf("Created user with username %s", username)
	log.Info().Msgf("Creating user with username %s", username)

	if username == "" {
		return nil, errors.Join(ErrMissingField, errors.New("username is required"))
	}
	if passwordHash == "" {
		return nil, errors.Join(ErrMissingField, errors.New("passwordHash is required"))
	}
	user := &User{Username: username, PasswordHash: passwordHash}
	createdID, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = createdID
	return user, nil
}

func (s *serviceImpl) GetByID(ctx context.Context, id int) (*User, error) {
	defer log.Info().Msgf("Got user by id %d", id)
	log.Info().Msgf("Getting user by id %d", id)
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) GetByUsername(ctx context.Context, username string) (*User, error) {
	defer log.Info().Msgf("Got user by username %s", username)
	log.Info().Msgf("Getting user by username %s", username)
	return s.repo.GetByUsername(ctx, username)
}

func (s *serviceImpl) Update(ctx context.Context, user *User) error {
	defer log.Info().Msgf("Updated user with id %d", user.ID)
	log.Info().Msgf("Updating user with id %d", user.ID)

	if user.Username == "" && user.PasswordHash == "" {
		return errors.Join(ErrMissingField, errors.New("username or passwordHash are required"))
	}

	return s.repo.Update(ctx, user)
}

func (s *serviceImpl) Delete(ctx context.Context, id int) error {
	defer log.Info().Msgf("Deleted user with id %d", id)
	log.Info().Msgf("Deleting user with id %d", id)
	return s.repo.Delete(ctx, id)
}
