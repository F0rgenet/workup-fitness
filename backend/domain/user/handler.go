package user

import (
	"encoding/json"
	"net/http"
	"time"

	"workup_fitness/middleware"
	"workup_fitness/pkg/httpx"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetPrivateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		httpx.Unauthorized(w, "Unauthorized")
		return
	}

	var req GetProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.GetByID(ctx, req.ID)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	var resp GetPrivateProfileResponse
	resp.ID = user.ID
	resp.Username = user.Username
	resp.CreatedAt = user.CreatedAt.Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	var req GetProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.GetByID(ctx, req.ID)

	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	var resp GetPublicProfileResponse
	resp.ID = user.ID
	resp.Username = user.Username
	resp.CreatedAt = user.CreatedAt.Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	var req UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	user := &User{
		ID:           req.ID,
		Username:     "",
		PasswordHash: "",
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			httpx.InternalServerError(w, err)
			return
		}

		user.PasswordHash = string(hashedPassword)
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	err := h.service.Update(ctx, user)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	var req DeleteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	err := h.service.Delete(ctx, req.ID)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
