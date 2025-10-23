package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guregu/null/v6/zero"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"workup_fitness/domain/auth/mocks"
	"workup_fitness/domain/user"
)

func TestRegisterHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService, "test-secret")

	expectedUser := &user.User{
		ID:        1,
		Username:  zero.StringFrom("testuser"),
		CreatedAt: time.Now(),
	}

	mockService.EXPECT().
		Register(gomock.Any(), "testuser", "password123").
		Return(expectedUser, nil)

	reqBody := RegisterRequest{
		Username: zero.StringFrom("testuser"),
		Password: zero.StringFrom("password123"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp AuthResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, "testuser", resp.User.Username.String)
	require.NotEmpty(t, resp.Token)

	token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService, "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRegisterHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService, "test-secret")

	mockService.EXPECT().
		Register(gomock.Any(), "testuser", "password123").
		Return(nil, errors.New("username already exists"))

	reqBody := RegisterRequest{
		Username: zero.StringFrom("testuser"),
		Password: zero.StringFrom("password123"),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRegisterHandler_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService, "test-secret")

	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	rr := httptest.NewRecorder()

	handler.Register(rr, req)
	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestLoginHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService, "test-secret")

	expectedUser := &user.User{
		ID:        1,
		Username:  zero.StringFrom("testuser"),
		CreatedAt: time.Now(),
	}

	mockService.EXPECT().
		Login(gomock.Any(), "testuser", "password123").
		Return(expectedUser, nil)

	reqBody := LoginRequest{
		Username: zero.StringFrom("testuser"),
		Password: zero.StringFrom("password123"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp AuthResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, "testuser", resp.User.Username.String)
	require.NotEmpty(t, resp.Token)

	token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService, "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService, "test-secret")

	mockService.EXPECT().
		Login(gomock.Any(), "testuser", "wrongpassword").
		Return(nil, ErrInvalidCreds)

	reqBody := LoginRequest{
		Username: zero.StringFrom("testuser"),
		Password: zero.StringFrom("wrongpassword"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestLoginHandler_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService, "test-secret")

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()

	handler.Login(rr, req)
	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPrepareAuthResponse_Success(t *testing.T) {
	testUser := &user.User{
		ID:        1,
		Username:  zero.StringFrom("testuser"),
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	resp, err := prepareAuthReponse(testUser, "test-secret")

	require.NoError(t, err)
	require.Equal(t, 1, resp.User.ID)
	require.Equal(t, "testuser", resp.User.Username.String)
	require.NotEmpty(t, resp.Token)

	token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	require.Equal(t, 1, int(claims["userID"].(float64)))
}
