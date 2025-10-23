package user

import (
	"context"

	"github.com/guregu/null/v6/zero"
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
	res := &serviceImpl{repo: repo}
	log.Info().Msg("Created user service")
	return res
}

func (s *serviceImpl) Create(ctx context.Context, username, passwordHash string) (*User, error) {
	log.Info().Msgf("Creating user with username %s", username)

	newName := zero.StringFromPtr(&username)
	newPasswordHash := zero.StringFromPtr(&passwordHash)

	user := &User{Username: newName, PasswordHash: newPasswordHash}
	createdID, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = createdID
	log.Info().Msgf("Created user with username %s", username)
	return user, nil
}

func (s *serviceImpl) GetByID(ctx context.Context, id int) (*User, error) {
	log.Info().Msgf("Getting user by id %d", id)
	user, err := s.repo.GetByID(ctx, id)
	log.Info().Msgf("Got user by id %d", id)
	return user, err
}

func (s *serviceImpl) GetByUsername(ctx context.Context, username string) (*User, error) {
	log.Info().Msgf("Getting user by username %s", username)
	user, err := s.repo.GetByUsername(ctx, username)
	log.Info().Msgf("Got user by username %s", username)
	return user, err
}

func (s *serviceImpl) Update(ctx context.Context, user *User) error {
	log.Info().Msgf("Updating user with id %d", user.ID)
	err := s.repo.Update(ctx, user)
	log.Info().Msgf("Updated user with id %d", user.ID)
	return err
}

func (s *serviceImpl) Delete(ctx context.Context, id int) error {
	log.Info().Msgf("Deleting user with id %d", id)
	err := s.repo.Delete(ctx, id)
	log.Info().Msgf("Deleted user with id %d", id)
	return err
}
