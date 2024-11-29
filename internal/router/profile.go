package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func ProfileRoutes(e *echo.Group) {
	h := handlers.NewProfileHandler()
	profile := e.Group("/profile", middlewares.AuthMiddleware)
	profile.PATCH("/", h.UpdateProfile)
}
