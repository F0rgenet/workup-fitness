package simulator_test

import (
	"context"
	"errors"
	"testing"
	"workup_fitness/domain/simulator"
	"workup_fitness/domain/simulator/mocks"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := simulator.NewService(repo)
	ctx := context.Background()

	repo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(1, nil)

	newSimulator, err := svc.Create(ctx, "Leg extension", "Some description", 0, 100, 10)
	require.NoError(t, err)

	require.Equal(t, 1, newSimulator.ID)
	require.Equal(t, "Leg extension", newSimulator.Name.String)
	require.Equal(t, "Some description", newSimulator.Description)
	require.Equal(t, 0.0, newSimulator.MinWeight)
	require.Equal(t, 100.0, newSimulator.MaxWeight)
	require.Equal(t, 10.0, newSimulator.WeightIncrement)
}

func TestService_Create_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := simulator.NewService(repo)
	ctx := context.Background()

	repo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(0, errors.New("some error"))

	_, err := svc.Create(ctx, "Leg extension", "Some description", 0, 100, 10)
	require.Error(t, err)
}
