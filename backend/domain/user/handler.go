package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"workup_fitness/middleware"
	"workup_fitness/pkg/httpx"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type ServiceInterface interface {
	GetByID(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func getContextUserID(ctx context.Context) (int, error) {
	val := ctx.Value(middleware.UserIDKey)
	userID, ok := val.(int)
	if !ok {
		return 0, errors.New("user id not found")
	}
	return userID, nil
}

func (h *Handler) GetPrivateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	userID, err := getContextUserID(ctx)
	if err != nil {
		httpx.Unauthorized(w, "Unauthorized")
		return
	}

	user, err := h.service.GetByID(ctx, userID)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	var resp GetPrivateProfileResponse
	resp.ID = user.ID
	resp.Username = user.Username
	resp.CreatedAt = user.CreatedAt.Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpx.InternalServerError(w, err)
		return
	}
}

func (h *Handler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpx.BadRequest(w, "Invalid profile id")
		return
	}

	user, err := h.service.GetByID(ctx, userID)

	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	var resp GetPublicProfileResponse
	resp.ID = user.ID
	resp.Username = user.Username
	resp.CreatedAt = user.CreatedAt.Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpx.InternalServerError(w, err)
		return
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	userID, err := getContextUserID(ctx)
	if err != nil {
		httpx.Unauthorized(w, "Unauthorized")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	user := &User{
		ID:           userID,
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

	if err := h.service.Update(ctx, user); err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		httpx.InternalServerError(w, err)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httpx.MethodNotAllowed(w)
		return
	}

	ctx := r.Context()

	userID, err := getContextUserID(ctx)
	if err != nil {
		httpx.Unauthorized(w, "Unauthorized")
		return
	}

	err = h.service.Delete(ctx, userID)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
