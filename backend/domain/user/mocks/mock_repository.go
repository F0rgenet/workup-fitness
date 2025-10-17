package mocks

import (
	"context"
	"database/sql"
	"workup_fitness/domain/user"

	"github.com/mattn/go-sqlite3"
)

type MockRepository struct {
	users  map[int]*user.User
	nextID int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users:  make(map[int]*user.User),
		nextID: 1,
	}
}

func (r *MockRepository) Create(ctx context.Context, user *user.User) (int, error) {
	for _, existingUser := range r.users {
		if existingUser.Username == user.Username {
			return 0, sqlite3.Error{
				Code:         sqlite3.ErrConstraint,
				ExtendedCode: sqlite3.ErrConstraintUnique,
			}
		}
	}

	r.users[r.nextID] = user
	r.nextID++
	return r.nextID - 1, nil
}

func (r *MockRepository) GetByID(ctx context.Context, id int) (*user.User, error) {
	found, ok := r.users[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (r *MockRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *MockRepository) Update(ctx context.Context, user *user.User) error {
	found, ok := r.users[user.ID]
	if !ok {
		return sql.ErrNoRows
	}

	for id, existingUser := range r.users {
		if id != user.ID && existingUser.Username == user.Username {
			return sqlite3.Error{
				Code:         sqlite3.ErrConstraint,
				ExtendedCode: sqlite3.ErrConstraintUnique,
			}
		}
	}

	found.Username = user.Username
	found.PasswordHash = user.PasswordHash
	return nil
}

func (r *MockRepository) Delete(ctx context.Context, id int) error {
	_, ok := r.users[id]
	if !ok {
		return sql.ErrNoRows
	}
	delete(r.users, id)
	return nil
}
