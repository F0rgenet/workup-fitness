package simulator

import (
	"context"
	"database/sql"

	"workup_fitness/internal/dbutil"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks workup_fitness/domain/simulator Repository

type Repository interface {
	Create(ctx context.Context, simulator *Simulator) (int, error)
	GetByID(ctx context.Context, id int) (*Simulator, error)
	GetByName(ctx context.Context, name string) (*Simulator, error)
	Update(ctxt context.Context, simulator *Simulator) error
	Delete(ctx context.Context, id int) error
}

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (repo *sqliteRepository) Create(ctx context.Context, simulator *Simulator) (int, error) {
	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO simulators (name, description, min_weight, max_weight, weight_increment) VALUES (?, ?, ?, ?, ?)`,
		simulator.Name, simulator.Description, simulator.MinWeight, simulator.MaxWeight, simulator.WeightIncrement,
	)
	log.Info().Msgf("Created simulator with name %s", simulator.Name.String)
	if err := dbutil.ProcessInsertError(err, ErrAlreadyExists, ErrMissingField); err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (repo *sqliteRepository) GetByID(ctx context.Context, id int) (*Simulator, error) {
	var simulator Simulator
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, description, min_weight, max_weight, weight_increment, created_at FROM simulators WHERE id = ?`,
		id,
	)
	err := row.Scan(&simulator.ID, &simulator.Name, &simulator.Description, &simulator.MinWeight, &simulator.MaxWeight, &simulator.WeightIncrement, &simulator.CreatedAt)
	if err != nil {
		return nil, err
	}
	if err := dbutil.ProcessRowError(err, ErrSimulatorNotFound); err != nil {
		return nil, err
	}
	return &simulator, err
}

func (repo *sqliteRepository) GetByName(ctx context.Context, name string) (*Simulator, error) {
	var simulator Simulator
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, description, min_weight, max_weight, weight_increment, created_at FROM simulators WHERE name = ?`,
		name,
	)
	err := row.Scan(&simulator.ID, &simulator.Name, &simulator.Description, &simulator.MinWeight, &simulator.MaxWeight, &simulator.WeightIncrement, &simulator.CreatedAt)
	dbutil.ProcessRowError(err, ErrSimulatorNotFound)
	return &simulator, nil
}

func (repo *sqliteRepository) Update(ctx context.Context, simulator *Simulator) error {
	_, err := repo.db.ExecContext(ctx,
		`UPDATE simulators SET name = ?, description = ?, min_weight = ?, max_weight = ?, weight_increment = ? WHERE id = ?`,
		simulator.Name, simulator.Description, simulator.MinWeight, simulator.MaxWeight, simulator.WeightIncrement, simulator.ID,
	)
	if err := dbutil.ProcessInsertError(err, ErrAlreadyExists, ErrMissingField); err != nil {
		return err
	}
	return err
}

func (repo *sqliteRepository) Delete(ctx context.Context, id int) error {
	result, err := repo.db.ExecContext(ctx,
		`DELETE FROM simulators WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrSimulatorNotFound
	}
	return nil
}
