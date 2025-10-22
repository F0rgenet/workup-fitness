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
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Creating user handler...")
	defer log.Info().Msg("Created user handler")
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

	log.Info().Msgf("Getting private profile for user with id %d", userID)

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

	log.Info().Msgf("Got private profile for user with id %d", userID)
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

	log.Info().Msgf("Getting public profile for user with id %d", userID)

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

	log.Info().Msgf("Got public profile for user with id %d", userID)
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

	log.Info().Msgf("Updating user with id %d", userID)

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

	// TODO: Вынести в ChangePassword в auth, заменить на UpdateUsername

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
		if errors.Is(err, ErrMissingField) {
			httpx.BadRequest(w, "Nothing to update, provide at least one field")
			return
		}

		httpx.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	log.Info().Msgf("Updated user with id %d", userID)
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

	log.Info().Msgf("Deleting user with id %d", userID)

	err = h.service.Delete(ctx, userID)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	log.Info().Msgf("Deleted user with id %d", userID)
}
