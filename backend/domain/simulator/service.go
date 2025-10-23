package simulator

import (
	"context"
	"fmt"
	"math"

	"github.com/guregu/null/v6/zero"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks workup_fitness/domain/simulator Service

type Service interface {
	Create(ctx context.Context, name, description string, minWeight, maxWeight, weightIncrement float64) (*Simulator, error)
	GetByID(ctx context.Context, id int) (*Simulator, error)
	GetByName(ctx context.Context, name string) (*Simulator, error)
	Update(ctx context.Context, simulator *Simulator) error
	Delete(ctx context.Context, id int) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) *serviceImpl {
	log.Info().Msg("Creating simulator service...")
	res := &serviceImpl{repo: repo}
	log.Info().Msg("Created simulator service")
	return res
}

func weightCheck(minWeight, maxWeight, weightIncrement float64) error {
	if minWeight < 0 || maxWeight < 0 {
		return ErrNegativeWeight
	}

	if weightIncrement == 0 {
		return ErrZeroIncrement
	}

	if minWeight == maxWeight {
		return fmt.Errorf("minWeight (%f) and maxWeight (%f) cannot be equal", minWeight, maxWeight)
	}

	if math.Abs(maxWeight-minWeight) < weightIncrement {
		return fmt.Errorf("maxWeight (%f) and minWeight (%f) are too close, increment (%f) is too small", maxWeight, minWeight, weightIncrement)
	}

	if minWeight < maxWeight && weightIncrement < 0 {
		return fmt.Errorf("weightIncrement cannot be negative when minWeight (%f) is less than maxWeight (%f)", minWeight, maxWeight)
	}

	if minWeight > maxWeight && weightIncrement > 0 {
		return fmt.Errorf("weightIncrement cannot be positive when minWeight (%f) is greater than maxWeight (%f)", minWeight, maxWeight)
	}
	return nil
}

func (s *serviceImpl) Create(ctx context.Context, name, description string, minWeight, maxWeight, weightIncrement float64) (*Simulator, error) {
	err := weightCheck(minWeight, maxWeight, weightIncrement)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Creating simulator with name %s", name)

	newName := zero.StringFromPtr(&name)

	simulator := &Simulator{Name: newName, Description: description, MinWeight: minWeight, MaxWeight: maxWeight, WeightIncrement: weightIncrement}
	createdID, err := s.repo.Create(ctx, simulator)
	if err != nil {
		return nil, err
	}
	simulator.ID = createdID
	log.Info().Msgf("Created simulator with name %s", simulator.Name.String)
	return simulator, nil
}

func (s *serviceImpl) GetByID(ctx context.Context, id int) (*Simulator, error) {
	log.Info().Msgf("Getting simulator by id %d", id)
	res, err := s.repo.GetByID(ctx, id)
	log.Info().Msgf("Got simulator by id %d", id)
	return res, err
}

func (s *serviceImpl) GetByName(ctx context.Context, name string) (*Simulator, error) {
	log.Info().Msgf("Getting simulator by name %s", name)
	res, err := s.repo.GetByName(ctx, name)
	log.Info().Msgf("Got simulator by name %s", name)
	return res, err
}

func (s *serviceImpl) Update(ctx context.Context, simulator *Simulator) error {
	log.Info().Msgf("Updating simulator with id %d", simulator.ID)
	if err := weightCheck(simulator.MinWeight, simulator.MaxWeight, simulator.WeightIncrement); err != nil {
		return err
	}
	err := s.repo.Update(ctx, simulator)
	log.Info().Msgf("Updated simulator with id %d", simulator.ID)
	return err
}

func (s *serviceImpl) Delete(ctx context.Context, id int) error {
	log.Info().Msgf("Deleting simulator with id %d", id)
	err := s.repo.Delete(ctx, id)
	log.Info().Msgf("Deleted simulator with id %d", id)
	return err
}
