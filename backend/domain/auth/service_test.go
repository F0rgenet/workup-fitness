package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"workup_fitness/domain/user"
	"workup_fitness/domain/user/mocks"
)

func TestRegister_Success(t *testing.T) {
	mockUserService := &mocks.MockService{
		CreateFunc: func(ctx context.Context, username, passwordHash string) (*user.User, error) {
			return &user.User{
				ID:           1,
				Username:     username,
				PasswordHash: passwordHash,
				CreatedAt:    time.Now(),
			}, nil
		},
	}

	authService := NewService(mockUserService)

	result, err := authService.Register(context.Background(), "testuser", "password123")

	require.NoError(t, err)
	require.Equal(t, "testuser", result.Username)

	err = bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte("password123"))
	require.NoError(t, err)
}

func TestRegister_EmptyUsername(t *testing.T) {
	mockUserService := &mocks.MockService{}
	authService := NewService(mockUserService)

	result, err := authService.Register(context.Background(), "", "password123")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrMissingField)
}

func TestRegister_EmptyPassword(t *testing.T) {
	mockUserService := &mocks.MockService{}
	authService := NewService(mockUserService)

	result, err := authService.Register(context.Background(), "testuser", "")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrMissingField)
}

func TestRegister_UserServiceError(t *testing.T) {
	mockUserService := &mocks.MockService{
		CreateFunc: func(ctx context.Context, username, passwordHash string) (*user.User, error) {
			return nil, user.ErrAlreadyExists
		},
	}

	authService := NewService(mockUserService)

	result, err := authService.Register(context.Background(), "testuser", "password123")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, user.ErrAlreadyExists)
}

func TestLogin_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockUserService := &mocks.MockService{
		GetByUsernameFunc: func(ctx context.Context, username string) (*user.User, error) {
			return &user.User{
				ID:           1,
				Username:     "testuser",
				PasswordHash: string(hashedPassword),
				CreatedAt:    time.Now(),
			}, nil
		},
	}

	authService := NewService(mockUserService)

	result, err := authService.Login(context.Background(), "testuser", "password123")

	require.NoError(t, err)
	require.Equal(t, "testuser", result.Username)
	require.Equal(t, 1, result.ID)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockUserService := &mocks.MockService{
		GetByUsernameFunc: func(ctx context.Context, username string) (*user.User, error) {
			return nil, user.ErrUserNotFound
		},
	}

	authService := NewService(mockUserService)

	result, err := authService.Login(context.Background(), "nonexistent", "password123")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrInvalidCreds)
}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockUserService := &mocks.MockService{
		GetByUsernameFunc: func(ctx context.Context, username string) (*user.User, error) {
			return &user.User{
				ID:           1,
				Username:     "testuser",
				PasswordHash: string(hashedPassword),
				CreatedAt:    time.Now(),
			}, nil
		},
	}

	authService := NewService(mockUserService)

	result, err := authService.Login(context.Background(), "testuser", "wrongpassword")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrInvalidCreds)
}

func TestLogin_ServiceError(t *testing.T) {
	mockUserService := &mocks.MockService{
		GetByUsernameFunc: func(ctx context.Context, username string) (*user.User, error) {
			return nil, errors.New("database connection error")
		},
	}

	authService := NewService(mockUserService)

	result, err := authService.Login(context.Background(), "testuser", "password123")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrInvalidCreds)
}
