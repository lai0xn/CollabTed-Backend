package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func WorkspaceRoutes(e *echo.Group) {
	h := handlers.NewWorkspaceHandler()
	workspaces := e.Group("/workspaces", middlewares.AuthMiddleware)
	workspaces.GET("/workspace", h.GetWorkspace)
	workspaces.POST("/create", h.CreateWorkspace)
	workspaces.POST("/invite", h.CreateInvitation)
	workspaces.GET("/accept", h.AcceptInvitation)
	workspaces.GET("/decline", h.DeclineInvitation)
}
