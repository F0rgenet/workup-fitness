package user_test

import (
	"context"
	"testing"
	"workup_fitness/domain/user"
	"workup_fitness/domain/user/mocks"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := user.NewService(repo)
	ctx := context.Background()

	repo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(1, nil)

	newUser, err := svc.Create(ctx, "alice", "some_hash")
	require.NoError(t, err)
	require.Equal(t, "alice", newUser.Username)
	require.Equal(t, "some_hash", newUser.PasswordHash)
	require.Equal(t, 1, newUser.ID)
}

func TestService_Create_MissingFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := user.NewService(repo)
	ctx := context.Background()

	u, err := svc.Create(ctx, "", "hash")
	require.ErrorIs(t, err, user.ErrMissingField)
	require.Nil(t, u)

	u, err = svc.Create(ctx, "bob", "")
	require.ErrorIs(t, err, user.ErrMissingField)
	require.Nil(t, u)
}

func TestService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := user.NewService(repo)
	ctx := context.Background()

	expectedUser := &user.User{
		ID:           1,
		Username:     "bob",
		PasswordHash: "hash",
	}

	repo.EXPECT().
		GetByID(ctx, 1).
		Return(expectedUser, nil)

	found, err := svc.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, expectedUser.Username, found.Username)
	require.Equal(t, expectedUser.PasswordHash, found.PasswordHash)
}

func TestService_GetByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := user.NewService(repo)
	ctx := context.Background()

	expectedUser := &user.User{
		ID:           1,
		Username:     "bob",
		PasswordHash: "hash",
	}

	repo.EXPECT().
		GetByUsername(ctx, "bob").
		Return(expectedUser, nil)

	found, err := svc.GetByUsername(ctx, "bob")
	require.NoError(t, err)
	require.Equal(t, expectedUser.Username, found.Username)
	require.Equal(t, expectedUser.PasswordHash, found.PasswordHash)
}

func TestService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := user.NewService(repo)
	ctx := context.Background()

	updatedUser := &user.User{
		ID:           1,
		Username:     "alice",
		PasswordHash: "nothash",
	}

	repo.EXPECT().Update(ctx, updatedUser).Return(nil)
	err := svc.Update(ctx, updatedUser)
	require.NoError(t, err)

	require.Equal(t, updatedUser.Username, updatedUser.Username)
	require.Equal(t, updatedUser.PasswordHash, updatedUser.PasswordHash)
	require.Equal(t, updatedUser.ID, updatedUser.ID)
}

func TestService_Update_MissingFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := user.NewService(repo)
	ctx := context.Background()

	updatedUser := &user.User{
		ID: 1,
	}

	err := svc.Update(ctx, updatedUser)
	require.ErrorIs(t, err, user.ErrMissingField)
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := user.NewService(repo)
	ctx := context.Background()

	repo.EXPECT().
		Delete(ctx, 1).
		Return(nil)

	err := svc.Delete(ctx, 1)
	require.NoError(t, err)
}
