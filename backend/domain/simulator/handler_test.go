package simulator_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/guregu/null/v6/zero"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"workup_fitness/domain/simulator"
	"workup_fitness/domain/simulator/mocks"
)

func TestCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockSimulator := &simulator.Simulator{
		ID:              1,
		Name:            zero.StringFrom("Bench Press"),
		Description:     "Chest exercise",
		MinWeight:       20.0,
		MaxWeight:       200.0,
		WeightIncrement: 2.5,
	}

	mockService.EXPECT().
		Create(gomock.Any(), "Bench Press", "Chest exercise", 20.0, 200.0, 2.5).
		Return(mockSimulator, nil)

	reqBody := simulator.CreateRequest{
		Name:            "Bench Press",
		Description:     "Chest exercise",
		MinWeight:       20.0,
		MaxWeight:       200.0,
		WeightIncrement: 2.5,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/simulators", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var resp simulator.CreateResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, mockSimulator.Name.String, resp.Name)
	require.Equal(t, mockSimulator.MinWeight, resp.MinWeight)
}

func TestCreate_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/simulators", bytes.NewReader([]byte("invalid json")))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreate_NegativeWeightError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockService.EXPECT().
		Create(gomock.Any(), "Test", "Description", -10.0, 100.0, 2.5).
		Return(nil, simulator.ErrNegativeWeight)

	reqBody := simulator.CreateRequest{
		Name:            "Test",
		Description:     "Description",
		MinWeight:       -10.0,
		MaxWeight:       100.0,
		WeightIncrement: 2.5,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/simulators", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreate_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockService.EXPECT().
		Create(gomock.Any(), "Test", "Description", 10.0, 100.0, 2.5).
		Return(nil, errors.New("db error"))

	reqBody := simulator.CreateRequest{
		Name:            "Test",
		Description:     "Description",
		MinWeight:       10.0,
		MaxWeight:       100.0,
		WeightIncrement: 2.5,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/simulators", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockSimulator := &simulator.Simulator{
		ID:              1,
		Name:            zero.StringFrom("Squat"),
		Description:     "Leg exercise",
		MinWeight:       40.0,
		MaxWeight:       300.0,
		WeightIncrement: 5.0,
	}

	mockService.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(mockSimulator, nil)

	req := httptest.NewRequest(http.MethodGet, "/simulators/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.GetByID(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp simulator.GetByIDResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, mockSimulator.Name.String, resp.Name)
	require.Equal(t, mockSimulator.MinWeight, resp.MinWeight)
}

func TestGetByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/simulators/invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.GetByID(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetByID_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockService.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/simulators/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.GetByID(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockService.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	reqBody := simulator.UpdateRequest{
		Name:            zero.StringFrom("Updated Name"),
		Description:     "Updated description",
		MinWeight:       30.0,
		MaxWeight:       250.0,
		WeightIncrement: 5.0,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/simulators/1", bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdate_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	reqBody := simulator.UpdateRequest{
		Name:            zero.StringFrom("Updated"),
		Description:     "Description",
		MinWeight:       30.0,
		MaxWeight:       250.0,
		WeightIncrement: 5.0,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/simulators/invalid", bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdate_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/simulators/1", bytes.NewReader([]byte("invalid json")))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdate_NegativeWeightError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockService.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(simulator.ErrNegativeWeight)

	reqBody := simulator.UpdateRequest{
		Name:            zero.StringFrom("Test"),
		Description:     "Description",
		MinWeight:       -10.0,
		MaxWeight:       100.0,
		WeightIncrement: 2.5,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/simulators/1", bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdate_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockService.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	reqBody := simulator.UpdateRequest{
		Name:            zero.StringFrom("Test"),
		Description:     "Description",
		MinWeight:       10.0,
		MaxWeight:       100.0,
		WeightIncrement: 2.5,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/simulators/1", bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockService.EXPECT().
		Delete(gomock.Any(), 1).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/simulators/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.Delete(rr, req)

	require.Equal(t, http.StatusAccepted, rr.Code)
}

func TestDelete_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/simulators/invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.Delete(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDelete_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := simulator.NewHandler(mockService)

	mockService.EXPECT().
		Delete(gomock.Any(), 1).
		Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodDelete, "/simulators/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.Delete(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
