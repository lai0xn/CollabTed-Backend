package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func WorkspaceRoutes(e *echo.Group) {
	h := handlers.NewWorkspaceHandler()
	workspaces := e.Group("/workspaces", middlewares.AuthMiddleware)
	workspaces.GET("/", h.GetWorkspaces)
	workspaces.POST("/create", h.CreateWorkspace)
	workspaces.GET("/:workspaceId", h.GetWorkspaceById)
	workspaces.GET("/connected/:workspaceId", h.GetConnectedUsers)
	workspaces.POST("/invite", h.InviteUser)
	workspaces.GET("/accept", h.AcceptInvitation)
	workspaces.GET("/:workspaceId/users", h.GetAllUsersInWorkspace)
	workspaces.GET("/:workspaceId/:userId", h.GetUserInWorkspace)
	workspaces.GET("/:workspaceId/invitations", h.GetAllInvites)
	workspaces.DELETE("/:invitationId/delete", h.DeleteInvitation)
	workspaces.DELETE("/:workspaceId", h.DeleteWorkspace)
	workspaces.PATCH("/:workspaceId/name", h.ChangeName)
	workspaces.POST("/:workspaceId/owner", h.ChangeOwner)
	workspaces.POST("/:workspaceId/:userId/role", h.ChangeUserRole)
	workspaces.POST("/:workspaceId/:userId/kick", h.KickUser)

}
