package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CollabTED/CollabTed-Backend/pkg/redis"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/pkg/utils"
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

func (s *WorkspaceService) GetWorkspace(workspaceId string) (*db.WorkspaceModel, error) {
	result, err := prisma.Client.Workspace.FindUnique(
		db.Workspace.ID.Equals(workspaceId),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CreateInvitation(email, workspaceID string) (*types.InvitationD, error) {
	token, err := utils.GenerateInvitationToken()
	if err != nil {
		return nil, err
	}

	invitation := types.InvitationD{
		Email:       email,
		Token:       token,
		WorkspaceID: workspaceID,
		Status:      string(types.PENDING),
	}

	err = redis.GetClient().Set(context.Background(), token, invitation, 132*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &invitation, nil
}

func AcceptInvitation(token string) (*types.InvitationD, error) {
	invitationData, err := redis.GetClient().Get(context.Background(), token).Result()
	if err != nil {
		return nil, err
	}

	var invitation types.InvitationD
	err = json.Unmarshal([]byte(invitationData), &invitation)
	if err != nil {
		return nil, err
	}

	invitation.Status = string(types.ACCEPTED)

	err = redis.GetClient().Set(context.Background(), token, invitation, 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	redis.GetClient().Del(context.Background(), token)

	return &invitation, nil
}

func DeclineInvitation(token string) (*types.InvitationD, error) {
	invitationData, err := redis.GetClient().Get(context.Background(), token).Result()
	if err != nil {
		return nil, err
	}

	var invitation types.InvitationD
	err = json.Unmarshal([]byte(invitationData), &invitation)
	if err != nil {
		return nil, err
	}

	invitation.Status = string(types.DECLINED)

	err = redis.GetClient().Set(context.Background(), token, invitation, 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	redis.GetClient().Del(context.Background(), token)

	return &invitation, nil
}
