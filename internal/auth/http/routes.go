package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(router chi.Router, handler *AuthHandler) {
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
}
