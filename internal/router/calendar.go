package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func CalendarRoutes(e *echo.Group) {
	h := handlers.NewCalendarHandler()
	events := e.Group("/events", middlewares.AuthMiddleware)
	events.POST("/create", h.CreateEvent)
	events.DELETE("/:eventId", h.DeleteEvent)
	events.GET("/list/:workspaceId", h.ListEvents)
	events.PUT("/:eventId", h.EditEvent)
}
