package user

import (
	"workup_fitness/config"
	"workup_fitness/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Get("/users/{id}", h.GetPublicProfile)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(config.JwtSecret))
		r.Get("/me", h.GetPrivateProfile)
		r.Put("/profile/update", h.Update)
		r.Delete("/profile/delete", h.Delete)
	})
}
