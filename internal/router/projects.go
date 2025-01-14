package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func ProjectsRoutes(e *echo.Group) {
	h := handlers.NewProjectHandler()

	projects := e.Group("/projects", middlewares.AuthMiddleware)
	projects.POST("/", h.CreateProject)          // Create a new project
	projects.GET("/:workspaceID", h.GetProjects) // List projects in a workspace
	projects.PUT("/:projectID", h.UpdateProject)
	projects.DELETE("/:projectID", h.DeleteProject)
	projects.GET("/project/:projectId", h.GetProjectById) // Get project by ID
}
