package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"workup_fitness/domain/user"
	"workup_fitness/pkg/httpx"
)

type Handler struct {
	service Service
	secret  string
}

func NewHandler(service Service, secret string) *Handler {
	log.Info().Msg("Creating auth handler...")
	defer log.Info().Msg("Created auth handler")
	return &Handler{service: service, secret: secret}
}

func prepareAuthReponse(user *user.User, secret string) (AuthResponse, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return AuthResponse{}, err
	}
	resp := AuthResponse{
		Token: tokenString,
		User: &UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}
	return resp, nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.MethodNotAllowed(w)
		return
	}

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.Register(r.Context(), req.Username.String, req.Password.String)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	log.Info().Msgf("Registered user with username %s", req.Username.String)

	resp, err := prepareAuthReponse(user, h.secret)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	log.Info().Msgf("Registered user with username %s", req.Username.String)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.MethodNotAllowed(w)
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.Login(r.Context(), req.Username.String, req.Password.String)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	log.Info().Msgf("Logged in user with username %s", req.Username.String)

	resp, err := prepareAuthReponse(user, h.secret)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	log.Info().Msgf("Logged in user with username %s", req.Username.String)
}
