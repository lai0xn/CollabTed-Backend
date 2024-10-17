package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func TasksRoutes(e *echo.Group) {
	tasks := e.Group("/tasks", middlewares.AuthMiddleware)
	taskHandler := handlers.NewTaskHandler()
	tasks.POST("/tasks", taskHandler.CreateTaskHandler)
	tasks.GET("/tasks/:id", taskHandler.GetTaskByIdHandler)
	tasks.GET("/projects/:projectId/tasks", taskHandler.ListTasksByProjectHandler)
	tasks.POST("/tasks/:id/assignees", taskHandler.AddAssigneeToTaskHandler)
}
