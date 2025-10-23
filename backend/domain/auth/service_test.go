package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"workup_fitness/domain/user"
	"workup_fitness/domain/user/mocks"
)

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)

	mockUserService.EXPECT().
		Create(gomock.Any(), "testuser", gomock.Any()).
		DoAndReturn(func(ctx context.Context, username, passwordHash string) (*user.User, error) {
			return &user.User{
				ID:           1,
				Username:     zero.StringFrom(username),
				PasswordHash: zero.StringFrom(passwordHash),
				CreatedAt:    time.Now(),
			}, nil
		})

	authService := NewService(mockUserService)

	result, err := authService.Register(context.Background(), "testuser", "password123")

	require.NoError(t, err)
	require.Equal(t, "testuser", result.Username)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(result.PasswordHash.String), []byte("password123")))
}

func TestRegister_EmptyUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	authService := NewService(mockUserService)

	result, err := authService.Register(context.Background(), "", "password123")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrMissingField)
}

func TestRegister_EmptyPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)
	authService := NewService(mockUserService)

	result, err := authService.Register(context.Background(), "testuser", "")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrMissingField)
}

func TestRegister_UserServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)

	mockUserService.EXPECT().
		Create(gomock.Any(), "testuser", gomock.Any()).
		Return(nil, user.ErrAlreadyExists)

	authService := NewService(mockUserService)

	result, err := authService.Register(context.Background(), "testuser", "password123")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, user.ErrAlreadyExists)
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	expectedUser := &user.User{
		ID:           1,
		Username:     zero.StringFrom("testuser"),
		PasswordHash: zero.StringFrom(string(hashedPassword)),
		CreatedAt:    time.Now(),
	}

	mockUserService.EXPECT().
		GetByUsername(gomock.Any(), "testuser").
		Return(expectedUser, nil)

	authService := NewService(mockUserService)

	result, err := authService.Login(context.Background(), "testuser", "password123")

	require.NoError(t, err)
	require.Equal(t, expectedUser.Username, result.Username)
	require.Equal(t, expectedUser.ID, result.ID)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)

	mockUserService.EXPECT().
		GetByUsername(gomock.Any(), "nonexistent").
		Return(nil, user.ErrUserNotFound)

	authService := NewService(mockUserService)

	result, err := authService.Login(context.Background(), "nonexistent", "password123")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrInvalidCreds)
}

func TestLogin_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	mockUserService.EXPECT().
		GetByUsername(gomock.Any(), "testuser").
		Return(&user.User{
			ID:           1,
			Username:     zero.StringFrom("testuser"),
			PasswordHash: zero.StringFrom(string(hashedPassword)),
			CreatedAt:    time.Now(),
		}, nil)

	authService := NewService(mockUserService)

	result, err := authService.Login(context.Background(), "testuser", "wrongpassword")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrInvalidCreds)
}

func TestLogin_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockService(ctrl)

	mockUserService.EXPECT().
		GetByUsername(gomock.Any(), "testuser").
		Return(nil, errors.New("database connection error"))

	authService := NewService(mockUserService)

	result, err := authService.Login(context.Background(), "testuser", "password123")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrInvalidCreds)
}
