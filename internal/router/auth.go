package router

import (
	"github.com/CollabTed/CollabTed-Backend/internal/handlers"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Group) {
	h := handlers.NewAuthHandler()
	auth := e.Group("/auth")
	auth.POST("/register", h.Register)
	auth.GET("/verify", h.VerifyUser)
	auth.POST("/login", h.Login)

}
