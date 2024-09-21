package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	"github.com/labstack/echo/v4"
)

func WorkspaceRoutes(e *echo.Group) {
	h := handlers.NewWorkspaceHandler()
	workspaces := e.Group("/workspaces")
	workspaces.GET("/workspace", h.GetWorkspace)
	workspaces.POST("/create", h.CreateWorkspace)
	workspaces.GET("/invite", h.CreateInvitation)
	workspaces.GET("/accept", h.AcceptInvitation)
	workspaces.GET("/decline", h.DeclineInvitation)
}
