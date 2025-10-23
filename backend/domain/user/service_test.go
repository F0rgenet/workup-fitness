package user_test

import (
	"context"
	"testing"
	"workup_fitness/domain/user"
	"workup_fitness/domain/user/mocks"

	"github.com/guregu/null/v6/zero"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
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
	require.Equal(t, "alice", newUser.Username.String)
	require.Equal(t, "some_hash", newUser.PasswordHash.String)
	require.Equal(t, 1, newUser.ID)
}

func TestService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	svc := user.NewService(repo)
	ctx := context.Background()

	expectedUser := &user.User{
		ID:           1,
		Username:     zero.StringFrom("bob"),
		PasswordHash: zero.StringFrom("hash"),
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
		Username:     zero.StringFrom("bob"),
		PasswordHash: zero.StringFrom("hash"),
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
		Username:     zero.StringFrom("alice"),
		PasswordHash: zero.StringFrom("nothash"),
	}

	repo.EXPECT().Update(ctx, updatedUser).Return(nil)
	err := svc.Update(ctx, updatedUser)
	require.NoError(t, err)

	require.Equal(t, updatedUser.Username, updatedUser.Username)
	require.Equal(t, updatedUser.PasswordHash, updatedUser.PasswordHash)
	require.Equal(t, updatedUser.ID, updatedUser.ID)
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
