package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

// TaskHandler struct with services as dependencies
type TaskHandler struct {
	TaskService    *services.TaskService
	ProjectService *services.ProjectService
}

// NewTaskHandler creates a new TaskHandler instance
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		TaskService:    services.NewTaskService(),
		ProjectService: services.NewProjectService(),
	}
}

// CreateTaskHandler handles task creation
func (h *TaskHandler) CreateTaskHandler(c echo.Context) error {
	var taskData types.TaskD

	// Parse input
	if err := c.Bind(&taskData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userId := c.Get("userId").(string) // Assume userId is extracted from JWT middleware
	workspaceId := taskData.WorkspaceID

	// Check if user has permissions to create the task (manager or project lead)
	canPerform, err := h.TaskService.CanUserPerformAction(userId, workspaceId, taskData.ProjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if !canPerform {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You do not have permission to create tasks in this workspace"})
	}

	// Create the task
	task, err := h.TaskService.CreateTask(taskData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, task)
}

// GetTaskByIdHandler retrieves a task by ID, ensuring the user is a member of the project.
func (h *TaskHandler) GetTaskByIdHandler(c echo.Context) error {
	taskID := c.Param("id")
	userId := c.Get("userId").(string) // Assume userId is extracted from JWT middleware
	workspaceId := c.QueryParam("workspaceId")

	// Ensure the user is a member of the project
	isMember, err := h.ProjectService.IsUserMemberOfProject(userId, workspaceId, taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if !isMember {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not a member of this project"})
	}

	// Retrieve the task
	task, err := h.TaskService.GetTaskById(taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, task)
}

// ListTasksByProjectHandler lists tasks in a project, ensuring the user is a member of the project.
func (h *TaskHandler) ListTasksByProjectHandler(c echo.Context) error {
	projectID := c.Param("projectId")
	userId := c.Get("userId").(string) // Assume userId is extracted from JWT middleware
	workspaceId := c.QueryParam("workspaceId")

	// Ensure the user is a member of the project
	isMember, err := h.ProjectService.IsUserMemberOfProject(userId, workspaceId, projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if !isMember {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not a member of this project"})
	}

	// List tasks in the project
	tasks, err := h.TaskService.ListTasksByProject(projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tasks)
}

// AddAssigneeToTaskHandler adds an assignee to a task
func (h *TaskHandler) AddAssigneeToTaskHandler(c echo.Context) error {
	var requestData struct {
		UserID string `json:"userId" binding:"required"`
	}

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	taskID := c.Param("id")
	userId := c.Get("userId").(string) // Assume userId is extracted from JWT middleware
	workspaceId := c.QueryParam("workspaceId")

	// Check if user has permissions to add assignees (manager or project lead)
	canPerform, err := h.TaskService.CanUserPerformAction(userId, workspaceId, taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if !canPerform {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You do not have permission to add assignees to this task"})
	}

	// Add the assignee to the task
	userWorkspace, err := h.TaskService.AddAssignee(workspaceId, taskID, requestData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, userWorkspace)
}
