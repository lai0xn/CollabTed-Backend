package services

import (
	"context"
	"fmt"

	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
)

// TaskService handles the task operations
type TaskService struct{}

// NewTaskService creates a new TaskService instance
func NewTaskService() *TaskService {
	return &TaskService{}
}

// CreateTask creates a new task in a project and assigns assignees.
func (s *TaskService) CreateTask(data types.TaskD) (*db.TaskModel, error) {
	// Create a new task
	result, err := prisma.Client.Task.CreateOne(
		db.Task.Project.Link(
			db.Project.ID.Equals(data.ProjectID),
		),
		db.Task.Title.Set(data.Title),
		db.Task.Description.Set(data.Description),
		db.Task.DueDate.Set(data.DueDate),
		db.Task.Priority.Set(db.Priority(data.Priority)),
		db.Task.Status.Link(
			db.Status.ID.Equals(data.StatusID),
		),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	// Link assignees to the task
	for _, assigneeID := range data.AssigneesIDs {
		_, err := prisma.Client.UserWorkspace.FindMany(
			db.UserWorkspace.UserID.Equals(assigneeID),
			db.UserWorkspace.WorkspaceID.Equals(data.WorkspaceID),
		).Update(
			db.UserWorkspace.Tasks.Link(
				db.Task.ID.Equals(result.ID), // Link the created task
			),
		).Exec(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to add assignee with ID %s to the task: %v", assigneeID, err)
		}
	}

	return result, nil
}

// GetTaskById retrieves a task by its ID.
func (s *TaskService) GetTaskById(taskID string) (*db.TaskModel, error) {
	task, err := prisma.Client.Task.FindUnique(
		db.Task.ID.Equals(taskID),
	).With(db.Task.Assignees.Fetch(), db.Task.Project.Fetch(), db.Task.Status.Fetch()).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return task, nil
}

// ListTasksByProject lists all tasks in a project.
func (s *TaskService) ListTasksByProject(projectID string) ([]db.TaskModel, error) {
	tasks, err := prisma.Client.Task.FindMany(
		db.Task.ProjectID.Equals(projectID),
	).With(db.Task.Assignees.Fetch(), db.Task.Status.Fetch()).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// AddAssignee adds a user to a task as an assignee.
func (s *TaskService) AddAssignee(workspaceID, taskID, userID string) (*db.UserWorkspaceModel, error) {
	ctx := context.Background()
	user, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userID),
		db.UserWorkspace.WorkspaceID.Equals(workspaceID),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	_, err = prisma.Client.Task.FindUnique(
		db.Task.ID.Equals(taskID),
	).Update(
		db.Task.Assignees.Link(db.UserWorkspace.ID.Equals(user.ID)),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CanUserPerformAction checks if a user has the required role to perform an action on a task.
func (s *TaskService) CanUserPerformAction(userId, workspaceId, taskId string) (bool, error) {
	// Find the UserWorkspace entry for the user in the workspace
	userWorkspace, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userId),
		db.UserWorkspace.WorkspaceID.Equals(workspaceId),
	).Exec(context.Background())

	if err != nil {
		return false, err
	}

	if userWorkspace == nil {
		return false, fmt.Errorf("user is not part of the workspace")
	}

	// Check if the user is a manager
	if userWorkspace.Role == db.UserRoleManager {
		return true, nil
	}

	// Fetch the task and check if the user is the lead of the associated project
	task, err := prisma.Client.Task.FindUnique(
		db.Task.ID.Equals(taskId),
	).With(db.Task.Project.Fetch()).Exec(context.Background())

	if err != nil {
		return false, err
	}

	// Check if the user is the lead of the project associated with the task
	if task.Project().LeadID == userWorkspace.ID {
		return true, nil
	}

	return false, nil
}

// AssignUserToTask assigns a single user to a task using the userWorkspaceID.
func (s *TaskService) AssignUserToTask(taskID, userWorkspaceID string) (*db.TaskModel, error) {
	ctx := context.Background()

	// Find the task and link the user as an assignee using userWorkspaceID
	task, err := prisma.Client.Task.FindUnique(
		db.Task.ID.Equals(taskID),
	).Update(
		db.Task.Assignees.Link(db.UserWorkspace.ID.Equals(userWorkspaceID)),
	).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to assign userWorkspaceID %s to task: %v", userWorkspaceID, err)
	}

	return task, nil
}

func (s *ProjectService) IsUserMemberOfProject(userId, workspaceId, projectId string) (bool, error) {
	// Check if the user is part of the workspace and project
	userWorkspace, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userId),
		db.UserWorkspace.WorkspaceID.Equals(workspaceId),
	).With(db.UserWorkspace.Projects.Fetch()).Exec(context.Background())

	if err != nil {
		return false, err
	}

	if userWorkspace == nil {
		return false, nil
	}

	// Check if the user is part of the specific project
	for _, project := range userWorkspace.Projects() {
		if project.ID == projectId {
			return true, nil
		}
	}

	return false, nil
}
