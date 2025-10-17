package user_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"
	"workup_fitness/domain/user"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	goose.SetLogger(goose.NopLogger())

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	err = goose.SetDialect("sqlite3")
	require.NoError(t, err)

	migrationsDir := filepath.Join("..", "..", "migrations")
	err = goose.Up(db, migrationsDir)
	require.NoError(t, err)

	return db
}

func newTestRepository(t *testing.T) (user.Repository, *sql.DB, context.Context) {
	t.Helper()

	db := setupTestDB(t)
	repo := user.NewSQLiteRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}

func TestRepository_Create(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	user := &user.User{
		Username:     "alice",
		PasswordHash: "some_hash235",
		CreatedAt:    time.Now(),
	}

	id, err := repo.Create(ctx, user)
	require.NoError(t, err)
	require.Greater(t, id, 0)
}

func TestRepository_GetByID(t *testing.T) {
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

func TestRepository_GetByUsername(t *testing.T) {
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

func TestRepository_Update(t *testing.T) {
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

func TestRepository_Delete(t *testing.T) {
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
