package user_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"workup_fitness/domain/user"
	"workup_fitness/internal/testutil"
)

func newTestRepository(t *testing.T) (user.Repository, *sql.DB, context.Context) {
	t.Helper()

	db := testutil.SetupTestDB(t)
	repo := user.NewSQLiteRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}

func TestRepository_Create_Success(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	user := &user.User{
		Username:     "alice",
		PasswordHash: "some_hash235",
	}

	id, err := repo.Create(ctx, user)
	require.NoError(t, err)
	require.Greater(t, id, 0)
}

func TestRepository_Create_AlreadyExists(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newUser := &user.User{
		Username:     "bob",
		PasswordHash: "hash456",
	}

	_, err := repo.Create(ctx, newUser)
	require.NoError(t, err)

	_, err = repo.Create(ctx, newUser)
	require.ErrorIs(t, err, user.ErrAlreadyExists)
}

func TestRepository_GetByID_Success(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newUser := &user.User{
		Username:     "bob",
		PasswordHash: "hash456",
		CreatedAt:    time.Now(),
	}

	id, err := repo.Create(ctx, newUser)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, newUser.Username, found.Username)
	require.Equal(t, newUser.PasswordHash, found.PasswordHash)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	found, err := repo.GetByID(ctx, 1)
	require.ErrorIs(t, err, user.ErrUserNotFound)
	require.Nil(t, found)
}

func TestRepository_GetByUsername_Success(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newUser := &user.User{
		Username:     "bobs",
		PasswordHash: "hash456",
		CreatedAt:    time.Now(),
	}

	_, err := repo.Create(ctx, newUser)
	require.NoError(t, err)

	found, err := repo.GetByUsername(ctx, newUser.Username)
	require.NoError(t, err)
	require.Equal(t, newUser.Username, found.Username)
	require.Equal(t, newUser.PasswordHash, found.PasswordHash)
}

func TestRepository_GetByUsername_NotFound(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	found, err := repo.GetByUsername(ctx, "bob")
	require.ErrorIs(t, err, user.ErrUserNotFound)
	require.Nil(t, found)
}

func TestRepository_Update_Success(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newUser := &user.User{
		Username:     "bob",
		PasswordHash: "hash456",
		CreatedAt:    time.Now(),
	}

	id, err := repo.Create(ctx, newUser)
	require.NoError(t, err)

	updatedUser := &user.User{
		ID:           id,
		Username:     "alice",
		PasswordHash: "nothash456",
	}
	err = repo.Update(ctx, updatedUser)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, updatedUser.Username, found.Username)
	require.Equal(t, updatedUser.PasswordHash, found.PasswordHash)
}

func TestRepository_Update_NotFound(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	updatedUser := &user.User{
		ID:           1,
		Username:     "alice",
		PasswordHash: "nothash456",
	}

	err := repo.Update(ctx, updatedUser)
	require.ErrorIs(t, err, user.ErrUserNotFound)
}

func TestRepository_Update_AlreadyExists(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	oldUser := &user.User{
		Username:     "bob",
		PasswordHash: "nothash456",
	}

	_, err := repo.Create(ctx, oldUser)
	require.NoError(t, err)

	newUser := &user.User{
		Username:     "alice",
		PasswordHash: "hash456",
	}

	_, err = repo.Create(ctx, newUser)
	require.NoError(t, err)

	updatedUser := &user.User{
		ID:           1,
		Username:     "alice",
		PasswordHash: "nothash456",
	}

	err = repo.Update(ctx, updatedUser)
	require.ErrorIs(t, err, user.ErrAlreadyExists)
}

func TestRepository_Delete_Success(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newUser := &user.User{
		Username:     "bob",
		PasswordHash: "hash456",
		CreatedAt:    time.Now(),
	}

	id, err := repo.Create(ctx, newUser)
	require.NoError(t, err)

	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, id)
	require.Error(t, err)
	require.Nil(t, found)
}

func TestRepository_Delete_NotFound(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	err := repo.Delete(ctx, 1)
	require.ErrorIs(t, err, user.ErrUserNotFound)
}
