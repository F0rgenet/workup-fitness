package mocks

import (
	"context"
	"errors"

	"workup_fitness/domain/user"
)

type MockService struct {
	RegisterFunc func(ctx context.Context, username, password string) (*user.User, error)
	LoginFunc    func(ctx context.Context, username, password string) (*user.User, error)
}

func (m *MockService) Register(ctx context.Context, username, password string) (*user.User, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, username, password)
	}
	return nil, errors.New("not implemented")
}

func (m *MockService) Login(ctx context.Context, username, password string) (*user.User, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, username, password)
	}
	return nil, errors.New("not implemented")
}
