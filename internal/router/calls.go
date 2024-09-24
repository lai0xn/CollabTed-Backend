package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func CallsRoutes(e *echo.Group) {
	h := handlers.NewCallHandler()

	calls := e.Group("/calls", middlewares.AuthMiddleware)
	calls.GET("/join/global/:workspaceId/:participantName", h.GetGlobalJoinToken)

}
