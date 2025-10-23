package user

import (
	"context"
	"database/sql"
	"workup_fitness/internal/dbutil"
)

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks workup_fitness/domain/user Repository

type Repository interface {
	Create(ctx context.Context, user *User) (int, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
}

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (repo *sqliteRepository) Create(ctx context.Context, user *User) (int, error) {
	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO users (username, password_hash) VALUES (?, ?)`,
		user.Username, user.PasswordHash,
	)
	if err := dbutil.ProcessInsertError(err, ErrAlreadyExists, ErrMissingField); err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (repo *sqliteRepository) GetByID(ctx context.Context, id int) (*User, error) {
	var user User
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE id = ?`,
		id,
	)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err := dbutil.ProcessRowError(err, ErrUserNotFound); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *sqliteRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE username = ?`,
		username,
	)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err := dbutil.ProcessRowError(err, ErrUserNotFound); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *sqliteRepository) Update(ctx context.Context, user *User) error {
	result, err := repo.db.ExecContext(ctx,
		`UPDATE users SET username = ?, password_hash = ? WHERE id = ?`,
		user.Username, user.PasswordHash, user.ID,
	)
	if err := dbutil.ProcessInsertError(err, ErrAlreadyExists, ErrMissingField); err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		err = ErrUserNotFound
	}
	return err
}

func (repo *sqliteRepository) Delete(ctx context.Context, id int) error {
	result, err := repo.db.ExecContext(ctx,
		`DELETE FROM users WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		err = ErrUserNotFound
	}

	return err
}
