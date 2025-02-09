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
	tasks.GET("/:workspaceId/:projectId/tasks", taskHandler.ListTasksByProjectHandler)
	tasks.POST("/:id/assignees", taskHandler.AddAssigneeToTaskHandler)
	tasks.DELETE("/:id/assignees/:assigneeId", taskHandler.RemoveAssigneeFromTaskHandler)
	tasks.PATCH("/:taskId/description", taskHandler.UpdateDescription)
	tasks.PATCH("/:taskId/title", taskHandler.UpdateTaskTitle)
	tasks.PATCH("/:taskId/:statusId/status", taskHandler.ChangeTaskStatus)
	tasks.DELETE("/:taskId", taskHandler.DeleteTaskHandler)
}
