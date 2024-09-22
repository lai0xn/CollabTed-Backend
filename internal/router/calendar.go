package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	"github.com/labstack/echo/v4"
)

func CalendarRoutes(e *echo.Group) {
	h := handlers.NewCalendarHandler()
	events := e.Group("/events")
	events.POST("/create", h.CreateEvent)
	events.GET("/list/:workspaceId", h.ListEvents)
}
