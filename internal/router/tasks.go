package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func TasksRoutes(e *echo.Group) {
	tasks := e.Group("/tasks", middlewares.AuthMiddleware)
	taskHandler := handlers.NewTaskHandler()
	tasks.POST("/", taskHandler.CreateTaskHandler)
	tasks.GET("/:id", taskHandler.GetTaskByIdHandler)
	tasks.GET("/:projectId/tasks", taskHandler.ListTasksByProjectHandler)
	tasks.POST("/:id/assignees", taskHandler.AddAssigneeToTaskHandler)
}
