package simulator_test

import (
	"context"
	"database/sql"
	"testing"

	"workup_fitness/domain/simulator"
	"workup_fitness/internal/testutil"

	"github.com/guregu/null/v6/zero"
	"github.com/stretchr/testify/require"
)

func newTestRepository(t *testing.T) (simulator.Repository, *sql.DB, context.Context) {
	t.Helper()

	db := testutil.SetupTestDB(t)
	repo := simulator.NewSQLiteRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}

func TestRepository_Create_Success(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	simulator := &simulator.Simulator{
		Name:            zero.StringFrom("Alice"),
		Description:     "Some description",
		MinWeight:       100,
		MaxWeight:       200,
		WeightIncrement: 10,
	}

	id, err := repo.Create(ctx, simulator)
	require.NoError(t, err)
	require.Equal(t, id, 1)
	require.NotNil(t, simulator.CreatedAt)
}

func TestRepository_Create_AlreadyExists(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newSimulator := &simulator.Simulator{
		Name:            zero.StringFrom("Bob"),
		Description:     "Some description",
		MinWeight:       100,
		MaxWeight:       200,
		WeightIncrement: 10,
	}

	_, err := repo.Create(ctx, newSimulator)
	require.NoError(t, err)

	_, err = repo.Create(ctx, newSimulator)
	require.ErrorIs(t, err, simulator.ErrAlreadyExists)
}

func TestRepository_Create_MissingFields(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	_, err := repo.Create(ctx, &simulator.Simulator{})
	require.ErrorIs(t, err, simulator.ErrMissingField)
	require.ErrorContains(t, err, "name")
}

func TestRepository_GetByID(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newSimulator := &simulator.Simulator{
		Name:            zero.StringFrom("Bob"),
		Description:     "Some description",
		MinWeight:       100,
		MaxWeight:       200,
		WeightIncrement: 10,
	}

	id, err := repo.Create(ctx, newSimulator)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, newSimulator.Name, found.Name)
	require.Equal(t, newSimulator.Description, found.Description)
	require.Equal(t, newSimulator.MinWeight, found.MinWeight)
	require.Equal(t, newSimulator.MaxWeight, found.MaxWeight)
	require.Equal(t, newSimulator.WeightIncrement, found.WeightIncrement)
}

func TestRepository_GetByName(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newSimulator := &simulator.Simulator{
		Name:            zero.StringFrom("Bob"),
		Description:     "Some description",
		MinWeight:       100,
		MaxWeight:       200,
		WeightIncrement: 10,
	}

	_, err := repo.Create(ctx, newSimulator)
	require.NoError(t, err)

	found, err := repo.GetByName(ctx, newSimulator.Name.String)
	require.NoError(t, err)
	require.Equal(t, newSimulator.Name, found.Name)
	require.Equal(t, newSimulator.Description, found.Description)
	require.Equal(t, newSimulator.MinWeight, found.MinWeight)
	require.Equal(t, newSimulator.MaxWeight, found.MaxWeight)
	require.Equal(t, newSimulator.WeightIncrement, found.WeightIncrement)
}

func TestRepository_Update(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newSimulator := &simulator.Simulator{
		Name:            zero.StringFrom("Bob"),
		Description:     "Some description",
		MinWeight:       100,
		MaxWeight:       200,
		WeightIncrement: 10,
	}

	id, err := repo.Create(ctx, newSimulator)
	require.NoError(t, err)

	updatedSimulator := &simulator.Simulator{
		ID:              id,
		Name:            zero.StringFrom("Alice"),
		Description:     "Some description",
		MinWeight:       100,
		MaxWeight:       200,
		WeightIncrement: 10,
	}
	err = repo.Update(ctx, updatedSimulator)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, updatedSimulator.Name, found.Name)
	require.Equal(t, updatedSimulator.Description, found.Description)
	require.Equal(t, updatedSimulator.MinWeight, found.MinWeight)
	require.Equal(t, updatedSimulator.MaxWeight, found.MaxWeight)
	require.Equal(t, updatedSimulator.WeightIncrement, found.WeightIncrement)
}

func TestRepository_Delete(t *testing.T) {
	repo, db, ctx := newTestRepository(t)
	defer db.Close()

	newSimulator := &simulator.Simulator{
		Name:            zero.StringFrom("Bob"),
		Description:     "Some description",
		MinWeight:       100,
		MaxWeight:       200,
		WeightIncrement: 10,
	}

	id, err := repo.Create(ctx, newSimulator)
	require.NoError(t, err)

	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, id)
	require.Error(t, err)
	require.Nil(t, found)
}
