package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"workup_fitness/domain/auth/mocks"
	"workup_fitness/domain/user"
)

func TestRegisterHandler_Success(t *testing.T) {
	mockService := &mocks.MockService{
		RegisterFunc: func(ctx context.Context, username, password string) (*user.User, error) {
			return &user.User{
				ID:        1,
				Username:  username,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewHandler(mockService, "test-secret")

	reqBody := RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp AuthResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "testuser", resp.User.Username)
	require.NotEmpty(t, resp.Token)

	token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	userID := int(claims["userID"].(float64))
	require.Equal(t, 1, userID)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	mockService := &mocks.MockService{}
	handler := NewHandler(mockService, "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRegisterHandler_ServiceError(t *testing.T) {
	mockService := &mocks.MockService{
		RegisterFunc: func(ctx context.Context, username, password string) (*user.User, error) {
			return nil, errors.New("username already exists")
		},
	}

	handler := NewHandler(mockService, "test-secret")

	reqBody := RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRegisterHandler_MethodNotAllowed(t *testing.T) {
	mockService := &mocks.MockService{}
	handler := NewHandler(mockService, "test-secret")

	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestLoginHandler_Success(t *testing.T) {
	mockService := &mocks.MockService{
		LoginFunc: func(ctx context.Context, username, password string) (*user.User, error) {
			return &user.User{
				ID:        1,
				Username:  username,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewHandler(mockService, "test-secret")

	reqBody := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp AuthResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "testuser", resp.User.Username)
	require.NotEmpty(t, resp.Token)

	token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	mockService := &mocks.MockService{}
	handler := NewHandler(mockService, "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginHandler_ServiceError(t *testing.T) {
	mockService := &mocks.MockService{
		LoginFunc: func(ctx context.Context, username, password string) (*user.User, error) {
			return nil, ErrInvalidCreds
		},
	}

	handler := NewHandler(mockService, "test-secret")

	reqBody := LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestLoginHandler_MethodNotAllowed(t *testing.T) {
	mockService := &mocks.MockService{}
	handler := NewHandler(mockService, "test-secret")

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPrepareAuthResponse_Success(t *testing.T) {
	testUser := &user.User{
		ID:        1,
		Username:  "testuser",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	resp, err := prepareAuthReponse(testUser, "test-secret")

	require.NoError(t, err)
	require.Equal(t, 1, resp.User.ID)
	require.Equal(t, "testuser", resp.User.Username)
	require.NotEmpty(t, resp.Token)

	token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	userID := int(claims["userID"].(float64))
	require.Equal(t, 1, userID)
}
