package simulator

import (
	"encoding/json"
	"net/http"
	"strconv"

	"workup_fitness/pkg/httpx"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	log.Info().Msg("Creating simulator handler...")
	res := &Handler{service: service}
	log.Info().Msg("Created simulator handler")
	return res
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	log.Info().Msgf("Creating simulator with name %s", req.Name)

	simulator, err := h.service.Create(ctx, req.Name, req.Description, req.MinWeight, req.MaxWeight, req.WeightIncrement)
	if err != nil {
		if err == ErrNegativeWeight || err == ErrZeroIncrement {
			httpx.BadRequest(w, err.Error())
			return
		}
		httpx.InternalServerError(w, err)
		return
	}

	var resp CreateResponse
	resp.ID = simulator.ID
	resp.Name = simulator.Name.String
	resp.Description = simulator.Description
	resp.MinWeight = simulator.MinWeight
	resp.MaxWeight = simulator.MaxWeight
	resp.WeightIncrement = simulator.WeightIncrement

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	log.Info().Msgf("Created simulator with id %d", simulator.ID)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	simulatorID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpx.BadRequest(w, "Invalid simulator id")
		return
	}

	log.Info().Msgf("Getting simulator with id %d", simulatorID)

	simulator, err := h.service.GetByID(ctx, simulatorID)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	var resp GetByIDResponse
	resp.ID = simulator.ID
	resp.Name = simulator.Name.String
	resp.Description = simulator.Description
	resp.MinWeight = simulator.MinWeight
	resp.MaxWeight = simulator.MaxWeight
	resp.WeightIncrement = simulator.WeightIncrement

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	log.Info().Msgf("Got simulator with id %d", simulatorID)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	simulatorID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpx.BadRequest(w, "Invalid simulator id")
		return
	}

	log.Info().Msgf("Updating simulator with id %d", simulatorID)

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	simulator := &Simulator{
		ID:              simulatorID,
		Name:            req.Name,
		Description:     req.Description,
		MinWeight:       req.MinWeight,
		MaxWeight:       req.MaxWeight,
		WeightIncrement: req.WeightIncrement,
	}

	if err := h.service.Update(ctx, simulator); err != nil {
		if err == ErrNegativeWeight || err == ErrZeroIncrement {
			httpx.BadRequest(w, err.Error())
			return
		}
		httpx.InternalServerError(w, err)
		return
	}

	var resp UpdateResponse
	resp.ID = simulator.ID
	resp.Name = simulator.Name.String
	resp.Description = simulator.Description
	resp.MinWeight = simulator.MinWeight
	resp.MaxWeight = simulator.MaxWeight
	resp.WeightIncrement = simulator.WeightIncrement

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	log.Info().Msgf("Updated simulator with id %d", simulatorID)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	simulatorID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpx.BadRequest(w, "Invalid simulator id")
		return
	}

	log.Info().Msgf("Deleting simulator with id %d", simulatorID)

	err = h.service.Delete(ctx, simulatorID)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	log.Info().Msgf("Deleted simulator with id %d", simulatorID)
}
