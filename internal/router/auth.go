package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Group) {
	h := handlers.NewAuthHandler()
	auth := e.Group("/auth")
	auth.POST("/register", h.Register)
	auth.GET("/verify", h.VerifyUser)
	auth.POST("/login", h.Login)
	auth.GET("/check", h.CheckUser, middlewares.AuthMiddleware)

}
