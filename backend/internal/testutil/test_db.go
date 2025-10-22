package testutil

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) *sql.DB {
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
