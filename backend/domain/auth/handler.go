package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"workup_fitness/domain/user"
	"workup_fitness/pkg/httpx"
)

type ServiceInterface interface {
	Register(ctx context.Context, username, password string) (*user.User, error)
	Login(ctx context.Context, username, password string) (*user.User, error)
}

type Handler struct {
	service ServiceInterface
	secret  string
}

func NewHandler(service ServiceInterface, secret string) *Handler {
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

	user, err := h.service.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

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

	user, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

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
}
