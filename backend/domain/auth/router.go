package auth

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post("/users/register", h.Register)
	r.Post("/users/login", h.Login)
}
