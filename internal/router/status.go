package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func StatusRoutes(e *echo.Group) {
	h := handlers.NewStatusHandler()
	statuses := e.Group("/statuses", middlewares.AuthMiddleware)
	statuses.POST("/create", h.CreateStatus)
	statuses.GET("/list/:projectId", h.GetStatusesByProject)
	statuses.GET("/:statusId", h.GetStatusByID)
	statuses.DELETE("/:statusId", h.DeleteStatus)
}
