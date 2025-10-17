package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"workup_fitness/domain/user"
	"workup_fitness/domain/user/mocks"
	"workup_fitness/middleware"
)

func TestGetPrivateProfile_Success(t *testing.T) {
	mockService := &mocks.MockService{
		GetByIDFunc: func(ctx context.Context, id int) (*user.User, error) {
			return &user.User{
				ID:        1,
				Username:  "testuser",
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			}, nil
		},
	}

	handler := user.NewHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetPrivateProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp user.GetPrivateProfileResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "testuser", resp.Username)
}

func TestGetPrivateProfile_Unauthorized(t *testing.T) {
	mockService := &mocks.MockService{}
	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rr := httptest.NewRecorder()

	handler.GetPrivateProfile(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetPrivateProfile_ServiceError(t *testing.T) {
	mockService := &mocks.MockService{
		GetByIDFunc: func(ctx context.Context, id int) (*user.User, error) {
			return nil, errors.New("database error")
		},
	}

	handler := user.NewHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetPrivateProfile(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetPublicProfile_Success(t *testing.T) {
	mockService := &mocks.MockService{
		GetByIDFunc: func(ctx context.Context, id int) (*user.User, error) {
			return &user.User{
				ID:        1,
				Username:  "publicuser",
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			}, nil
		},
	}

	handler := user.NewHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.GetPublicProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp user.GetPublicProfileResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "publicuser", resp.Username)
}

func TestGetPublicProfile_InvalidID(t *testing.T) {
	mockService := &mocks.MockService{}
	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.GetPublicProfile(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetPublicProfile_ServiceError(t *testing.T) {
	mockService := &mocks.MockService{
		GetByIDFunc: func(ctx context.Context, id int) (*user.User, error) {
			return nil, errors.New("database error")
		},
	}

	handler := user.NewHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.GetPublicProfile(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdate_Success(t *testing.T) {
	mockService := &mocks.MockService{
		UpdateFunc: func(ctx context.Context, user *user.User) error {
			return nil
		},
	}

	handler := user.NewHandler(mockService)

	reqBody := user.UpdateRequest{
		Username: "newusername",
		Password: "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/profile/update", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdate_Unauthorized(t *testing.T) {
	mockService := &mocks.MockService{}
	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/profile/update", nil)
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdate_ServiceError(t *testing.T) {
	mockService := &mocks.MockService{
		UpdateFunc: func(ctx context.Context, user *user.User) error {
			return errors.New("database error")
		},
	}

	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/profile/update", strings.NewReader(`{}`))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDelete_Success(t *testing.T) {
	mockService := &mocks.MockService{
		DeleteFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	handler := user.NewHandler(mockService)
	req := httptest.NewRequest(http.MethodDelete, "/profile/delete", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	require.Equal(t, http.StatusAccepted, rr.Code)
}

func TestDelete_Unauthorized(t *testing.T) {
	mockService := &mocks.MockService{}
	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/profile/delete", nil)
	rr := httptest.NewRecorder()

	handler.Delete(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDelete_ServiceError(t *testing.T) {
	mockService := &mocks.MockService{
		DeleteFunc: func(ctx context.Context, id int) error {
			return errors.New("database error")
		},
	}

	handler := user.NewHandler(mockService)
	req := httptest.NewRequest(http.MethodDelete, "/profile/delete", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
