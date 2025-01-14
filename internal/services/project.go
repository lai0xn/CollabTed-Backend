package services

import (
	"context"
	"fmt"

	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
)

type ProjectService struct{}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// CreateProject creates a new project in a workspace and assigns the lead and assignees.
func (s *ProjectService) CreateProject(data types.ProjectD) (*db.ProjectModel, error) {
	// Create a new project
	result, err := prisma.Client.Project.CreateOne(
		db.Project.Title.Set(data.Title),
		db.Project.Workspace.Link(
			db.Workspace.ID.Equals(data.WorksapceID),
		),
		db.Project.Lead.Link(
			db.UserWorkspace.ID.Equals(data.LeadID),
		),
	).With(db.Project.Assignees.Fetch()).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println(result.ID)
	// Link assignees to the project
	for _, assigneeID := range data.AssigneesIDs {
		usr, err := prisma.Client.UserWorkspace.FindFirst(
			db.UserWorkspace.UserID.Equals(assigneeID),
			db.UserWorkspace.WorkspaceID.Equals(data.WorksapceID),
		).Exec(context.Background())
		if err != nil {
			return nil, err
		}
		_, err = prisma.Client.UserWorkspace.FindUnique(
			db.UserWorkspace.ID.Equals(usr.ID),
		).Update(
			db.UserWorkspace.Projects.Link(
				db.Project.ID.Equals(result.ID),
			),
		).Exec(context.Background())
		if err != nil {
			return nil, err
		}

	}
	return result, nil
}

// GetProjectById retrieves a project by its ID.
func (s *ProjectService) GetProjectById(projectID string) (*db.ProjectModel, error) {
	project, err := prisma.Client.Project.FindUnique(
		db.Project.ID.Equals(projectID),
	).With(db.Project.Lead.Fetch(), db.Project.Statuses.Fetch(), db.Project.Assignees.Fetch()).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return project, nil
}

// ListProjectsByWorkspace lists all projects in a workspace.
func (s *ProjectService) ListProjectsByWorkspace(userID, workspaceID string) ([]db.ProjectModel, error) {
	// Check if the user is part of the workspace
	userWorkspace, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userID),
		db.UserWorkspace.WorkspaceID.Equals(workspaceID),
	).Exec(context.Background())
	if err != nil {
		fmt.Println(userID, workspaceID)
		fmt.Println(err.Error())
		return nil, err
	}

	if userWorkspace == nil {
		return nil, fmt.Errorf("user is not a member of the workspace")
	}

	// Proceed to fetch projects if the user is a member
	projects, err := prisma.Client.Project.FindMany(
		db.Project.WorkspaceID.Equals(workspaceID),
	).With(db.Project.Assignees.Fetch()).Exec(context.Background())

	if err != nil {
		return nil, err
	}
	return projects, nil
}

// AddAssignee adds a user to a project as an assignee.
func (s *ProjectService) AddAssignee(workspaceID, projectID, userID string) (*db.UserWorkspaceModel, error) {
	ctx := context.Background()
	user, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userID),
		db.UserWorkspace.WorkspaceID.Equals(workspaceID),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	_, err = prisma.Client.Project.FindUnique(
		db.Project.ID.Equals(projectID),
	).Update(
		db.Project.Assignees.Link(db.UserWorkspace.ID.Equals(user.ID)),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ProjectService) CanUserPerformAction(userId, workspaceId string, requiredRole db.UserRole) (bool, error) {
	userWorkspace, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userId),
		db.UserWorkspace.WorkspaceID.Equals(workspaceId),
	).Exec(context.Background())

	if err != nil {
		return false, err
	}

	if userWorkspace == nil {
		return false, nil
	}

	logger.LogInfo().Msg(string(userWorkspace.Role))
	logger.LogInfo().Msg(string(requiredRole))
	return userWorkspace.Role == requiredRole, nil
}

func (s *ProjectService) UpdateProject(data types.ProjectD, projectId string) (*db.ProjectModel, error) {
	ctx := context.Background()
	result, err := prisma.Client.Project.FindUnique(
		db.Project.ID.Equals(projectId),
	).Update(
		db.Project.Title.Set(data.Title),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ProjectService) DeleteProject(projectId string) error {
	ctx := context.Background()
	_, err := prisma.Client.Project.FindUnique(
		db.Project.ID.Equals(projectId),
	).Delete().Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
