package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
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
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return 0, ErrAlreadyExists
			}
		}
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (repo *sqliteRepository) Update(ctx context.Context, user *User) error {
	result, err := repo.db.ExecContext(ctx,
		`UPDATE users SET username = ?, password_hash = ? WHERE id = ?`,
		user.Username, user.PasswordHash, user.ID,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return ErrAlreadyExists
			}
		}
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
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
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
