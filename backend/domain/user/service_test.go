package user_test

import (
	"context"
	"testing"
	"workup_fitness/domain/user"
	"workup_fitness/domain/user/mocks"

	"github.com/stretchr/testify/require"
)

func TestService_Create_Success(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	newUser, err := svc.Create(ctx, "alice", "some_hash")
	require.NoError(t, err)
	require.Equal(t, "alice", newUser.Username)
	require.Equal(t, "some_hash", newUser.PasswordHash)
	require.Greater(t, newUser.ID, 0)
}

func TestService_Create_MissingFields(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	u, err := svc.Create(ctx, "", "hash")
	require.ErrorIs(t, err, user.ErrMissingField)
	require.Nil(t, u)

	u, err = svc.Create(ctx, "bob", "")
	require.ErrorIs(t, err, user.ErrMissingField)
	require.Nil(t, u)
}

func TestService_Create_AlreadyExists(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, "bob", "hash")
	require.NoError(t, err)

	_, err = svc.Create(ctx, "bob", "hash")
	require.ErrorIs(t, err, user.ErrAlreadyExists)
}

func TestService_GetByID(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	newUser := &user.User{
		Username:     "bob",
		PasswordHash: "hash",
	}
	_, err := repo.Create(ctx, newUser)
	require.NoError(t, err)

	found, err := svc.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, newUser.Username, found.Username)
	require.Equal(t, newUser.PasswordHash, found.PasswordHash)
}

func TestService_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	found, err := svc.GetByID(ctx, 1)
	require.ErrorIs(t, err, user.ErrUserNotFound)
	require.Nil(t, found)
}

func TestService_GetByUsername(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	newUser := &user.User{
		Username:     "bob",
		PasswordHash: "hash",
	}
	_, err := repo.Create(ctx, newUser)
	require.NoError(t, err)

	found, err := svc.GetByUsername(ctx, newUser.Username)
	require.NoError(t, err)
	require.Equal(t, newUser.Username, found.Username)
	require.Equal(t, newUser.PasswordHash, found.PasswordHash)
}

func TestService_GetByUsername_NotFound(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	found, err := svc.GetByUsername(ctx, "bob")
	require.ErrorIs(t, err, user.ErrUserNotFound)
	require.Nil(t, found)
}

func TestService_Update(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	newUser, err := svc.Create(ctx, "bob", "hash")
	require.NoError(t, err)

	updatedUser := &user.User{
		ID:           newUser.ID,
		Username:     "alice",
		PasswordHash: "nothash",
	}
	err = svc.Update(ctx, updatedUser)
	require.NoError(t, err)

	found, err := svc.GetByID(ctx, newUser.ID)
	require.NoError(t, err)
	require.Equal(t, updatedUser.Username, found.Username)
	require.Equal(t, updatedUser.PasswordHash, found.PasswordHash)
}

func TestService_Update_NotFound(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	updatedUser := &user.User{
		ID:           1,
		Username:     "alice",
		PasswordHash: "nothash",
	}
	err := svc.Update(ctx, updatedUser)
	require.ErrorIs(t, err, user.ErrUserNotFound)
}

func TestService_Update_AlreadyExists(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	newUser, err := svc.Create(ctx, "bob", "hash")
	require.NoError(t, err)

	_, err = svc.Create(ctx, "alice", "hash")
	require.NoError(t, err)

	updatedUser := &user.User{
		ID:           newUser.ID,
		Username:     "alice",
		PasswordHash: "nothash",
	}
	err = svc.Update(ctx, updatedUser)
	require.ErrorIs(t, err, user.ErrAlreadyExists)
}

func TestService_Delete(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	newUser, err := svc.Create(ctx, "bob", "hash")
	require.NoError(t, err)

	err = svc.Delete(ctx, newUser.ID)
	require.NoError(t, err)

	found, err := svc.GetByID(ctx, newUser.ID)
	require.ErrorIs(t, err, user.ErrUserNotFound)
	require.Nil(t, found)
}

func TestService_Delete_NotFound(t *testing.T) {
	repo := mocks.NewMockRepository()
	svc := user.NewService(repo)
	ctx := context.Background()

	err := svc.Delete(ctx, 1)
	require.ErrorIs(t, err, user.ErrUserNotFound)
}
