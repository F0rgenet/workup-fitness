package user

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, user *User) (int, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
}

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (repo *sqliteRepository) Create(ctx context.Context, user *User) (int, error) {
	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)`,
		user.Username, user.PasswordHash, user.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (repo *sqliteRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE username = ?`,
		username,
	)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *sqliteRepository) GetByID(ctx context.Context, id int) (*User, error) {
	var user User
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE id = ?`,
		id,
	)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}