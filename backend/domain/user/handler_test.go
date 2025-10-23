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
	"github.com/guregu/null/v6/zero"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"workup_fitness/domain/user"
	"workup_fitness/domain/user/mocks"
	"workup_fitness/middleware"
)

func TestGetPrivateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	mockUser := &user.User{
		ID:        1,
		Username:  zero.StringFrom("testuser"),
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mockService.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(mockUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.GetPrivateProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp user.GetPrivateProfileResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, mockUser.Username, resp.Username)
}

func TestGetPrivateProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rr := httptest.NewRecorder()

	handler.GetPrivateProfile(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetPrivateProfile_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	mockService.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.GetPrivateProfile(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetPublicProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	mockUser := &user.User{
		ID:        1,
		Username:  zero.StringFrom("publicuser"),
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mockService.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(mockUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.GetPublicProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp user.GetPublicProfileResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, mockUser.Username, resp.Username)
}

func TestGetPublicProfile_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	mockService.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.GetPublicProfile(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	reqBody := user.UpdateRequest{
		Username: zero.StringFrom("newusername"),
		Password: zero.StringFrom("newpassword"),
	}
	body, _ := json.Marshal(reqBody)

	mockService.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/profile/update", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdate_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	mockService.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodPut, "/profile/update", strings.NewReader(`{}`))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	mockService.EXPECT().
		Delete(gomock.Any(), 1).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/profile/delete", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.Delete(rr, req)

	require.Equal(t, http.StatusAccepted, rr.Code)
}

func TestDelete_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := user.NewHandler(mockService)

	mockService.EXPECT().
		Delete(gomock.Any(), 1).
		Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodDelete, "/profile/delete", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.Delete(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
