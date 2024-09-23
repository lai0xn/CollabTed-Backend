package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CollabTED/CollabTed-Backend/config"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/mail"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/pkg/utils"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
)

type WorkspaceService struct {
	sender *mail.EmailVerifier
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		sender: mail.NewVerifier(),
	}
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

	logger.LogInfo().Msg(string(userWorkspace.Role))
	logger.LogInfo().Msg(string(requiredRole))
	return userWorkspace.Role == requiredRole, nil
}

func (s *WorkspaceService) SendInvitation(email, workspaceID string) error {
	token, err := utils.GenerateInvitationToken()
	if err != nil {
		return err
	}

	_, err = prisma.Client.Invitation.CreateOne(
		db.Invitation.Email.Set(email),
		db.Invitation.Token.Set(token),
		db.Invitation.Workspace.Link(
			db.Workspace.ID.Equals(workspaceID),
		),
		db.Invitation.Status.Set(db.InvitationStatusPending),
	).Exec(context.Background())
	if err != nil {
		return err
	}

	invitationLink := fmt.Sprintf("%s/workspaces/join?token=%s", config.HOST_URL, token)

	if err := s.sender.SendInvitationMail([]string{email}, invitationLink); err != nil {
		return err
	}

	return nil
}

func (s *WorkspaceService) AcceptInvitation(userID, token string) error {
	invitation, err := prisma.Client.Invitation.FindUnique(
		db.Invitation.Token.Equals(token),
	).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	if invitation.Status != db.InvitationStatusPending {
		return fmt.Errorf("invitation has already been accepted or declined")
	}

	_, err = prisma.Client.UserWorkspace.CreateOne(
		db.UserWorkspace.User.Link(
			db.User.ID.Equals(userID),
		),
		db.UserWorkspace.Workspace.Link(
			db.Workspace.ID.Equals(invitation.WorkspaceID),
		),
		db.UserWorkspace.Role.Set(db.UserRoleMember),
		db.UserWorkspace.JoinedAt.Set(time.Now()),
	).Exec(context.Background())

	if err != nil {
		return fmt.Errorf("failed to join workspace: %v", err)
	}

	_, err = prisma.Client.Invitation.FindUnique(
		db.Invitation.Token.Equals(token),
	).Update(
		db.Invitation.Status.Set(db.InvitationStatusAccepted),
	).Exec(context.Background())

	if err != nil {
		return fmt.Errorf("failed to update invitation status: %v", err)
	}

	return nil
}

func (s *WorkspaceService) GetAllUsersInWorkspace(workspaceId string) ([]db.UserModel, error) {
	userWorkspaces, err := prisma.Client.UserWorkspace.FindMany(
		db.UserWorkspace.WorkspaceID.Equals(workspaceId),
	).Exec(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to get user-workspace relations: %v", err)
	}

	var userIds []string
	for _, userWorkspace := range userWorkspaces {
		userIds = append(userIds, userWorkspace.UserID)
	}

	if len(userIds) == 0 {
		return []db.UserModel{}, nil
	}

	users, err := prisma.Client.User.FindMany(
		db.User.ID.In(userIds),
	).Exec(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	return users, nil
}

func (s *WorkspaceService) GetInvitations(workspaceId string) ([]db.InvitationModel, error) {
	invitations, err := prisma.Client.Invitation.FindMany(
		db.Invitation.WorkspaceID.Equals(workspaceId),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (s *WorkspaceService) DeleteInvitation(invitationId string) error {
	_, err := prisma.Client.Invitation.FindUnique(
		db.Invitation.ID.Equals(invitationId),
	).Delete().Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}
