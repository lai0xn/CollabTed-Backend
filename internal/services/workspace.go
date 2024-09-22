package services

import (
	"context"
	"time"

	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
)

type WorkspaceService struct{}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

func (s *WorkspaceService) CreateWorkspace(data types.WorkspaceD) (*db.WorkspaceModel, error) {
	result, err := prisma.Client.Workspace.CreateOne(
		db.Workspace.WorkspaceName.Set(data.Name),
		db.Workspace.Owner.Link(
			db.User.ID.Equals(data.OwnerID),
		),
	).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	_, err = prisma.Client.UserWorkspace.CreateOne(
		db.UserWorkspace.User.Link(
			db.User.ID.Equals(data.OwnerID),
		),
		db.UserWorkspace.Workspace.Link(
			db.Workspace.ID.Equals(result.ID),
		),
		db.UserWorkspace.Role.Set(db.UserRoleAdmin),
		db.UserWorkspace.JoinedAt.Set(time.Now()),
	).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WorkspaceService) ListWorkspaces(userID string) ([]db.WorkspaceModel, error) {
	result, err := prisma.Client.Workspace.FindMany(
		db.Workspace.Users.Some(
			db.UserWorkspace.UserID.Equals(userID),
		),
	).Exec(context.Background())

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WorkspaceService) GetWorkspaceById(workspaceId string) (*db.WorkspaceModel, error) {
	result, err := prisma.Client.Workspace.FindUnique(
		db.Workspace.ID.Equals(workspaceId),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WorkspaceService) CanUserPerformAction(userId, workspaceId string, requiredRole db.UserRole) (bool, error) {
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

	return userWorkspace.Role == requiredRole, nil
}
