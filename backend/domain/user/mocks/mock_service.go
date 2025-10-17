package mocks

import (
	"context"
	"errors"
	"workup_fitness/domain/user"
)

type MockService struct {
	CreateFunc        func(ctx context.Context, username, passwordHash string) (*user.User, error)
	GetByIDFunc       func(ctx context.Context, id int) (*user.User, error)
	GetByUsernameFunc func(ctx context.Context, username string) (*user.User, error)
	UpdateFunc        func(ctx context.Context, user *user.User) error
	DeleteFunc        func(ctx context.Context, id int) error
}

func (m *MockService) Create(ctx context.Context, username, passwordHash string) (*user.User, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, username, passwordHash)
	}
	return nil, errors.New("not implemented")
}

func (m *MockService) GetByID(ctx context.Context, id int) (*user.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockService) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	if m.GetByUsernameFunc != nil {
		return m.GetByUsernameFunc(ctx, username)
	}
	return nil, errors.New("not implemented")
}

func (m *MockService) Update(ctx context.Context, user *user.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return errors.New("not implemented")
}

func (m *MockService) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}
