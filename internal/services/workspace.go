package services

import (
	"context"
	"encoding/json"
	"errors"
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
	sender   *mail.EmailVerifier
	boardSrv *BoardService
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		sender:   mail.NewVerifier(),
		boardSrv: NewBoardService(),
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

	_, err = s.boardSrv.SaveBoard(types.BoardD{
		WorkspaceID: result.ID,
		Elements:    []json.RawMessage{},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WorkspaceService) ListWorkspaces(userID string) ([]db.WorkspaceModel, error) {
	result, err := prisma.Client.Workspace.FindMany(
		db.Workspace.Or(
			db.Workspace.OwnerID.Equals(userID),
			db.Workspace.Users.Some(
				db.UserWorkspace.UserID.Equals(userID),
			),
		),
	).Exec(context.Background())

	if err != nil {
		return nil, err
	}
	logger.LogInfo().Msgf("Found %d workspaces for user %s", len(result), userID)

	return result, nil
}

func (s *WorkspaceService) GetWorkspaceById(workspaceId string) (*db.WorkspaceModel, error) {
	result, err := prisma.Client.Workspace.FindUnique(
		db.Workspace.ID.Equals(workspaceId),
	).With(db.Workspace.Users.Fetch()).Exec(context.Background())
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
	_, err = prisma.Client.User.FindFirst(
		db.User.Email.Equals(email),
	).Exec(context.Background())

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

func (s *WorkspaceService) AcceptInvitation(userID, token string) (string, error) {
	invitation, err := prisma.Client.Invitation.FindUnique(
		db.Invitation.Token.Equals(token),
	).Exec(context.Background())
	if err != nil {
		return "", fmt.Errorf("invitation not found")
	}

	if invitation.Status != db.InvitationStatusPending {
		return "", fmt.Errorf("invitation has already been accepted or declined")
	}

	user, err := prisma.Client.User.FindUnique(
		db.User.ID.Equals(userID),
	).Exec(context.Background())
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	existingUsers, err := s.GetAllUsersInWorkspace(invitation.WorkspaceID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve users in workspace: %v", err)
	}

	uniqueName := utils.GenerateUniqueName(user.Name, existingUsers)

	if uniqueName != user.Name {
		_, err := prisma.Client.User.FindUnique(
			db.User.ID.Equals(userID),
		).Update(
			db.User.Name.Set(uniqueName),
		).Exec(context.Background())
		if err != nil {
			return "", fmt.Errorf("failed to update user's name: %v", err)
		}
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
		return "", fmt.Errorf("failed to join workspace: %v", err)
	}

	_, err = prisma.Client.Invitation.FindUnique(
		db.Invitation.Token.Equals(token),
	).Update(
		db.Invitation.Status.Set(db.InvitationStatusAccepted),
	).Exec(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to update invitation status: %v", err)
	}

	return invitation.WorkspaceID, nil
}

func (s *WorkspaceService) GetAllUsersInWorkspace(workspaceId string) ([]types.UserWorkspace, error) {
	// Fetch the user-workspace relationships for the given workspace ID
	userWorkspaces, err := prisma.Client.UserWorkspace.FindMany(
		db.UserWorkspace.WorkspaceID.Equals(workspaceId),
	).Exec(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to get user-workspace relations: %v", err)
	}

	var userIds []string
	userWorkspaceMap := make(map[string]db.UserWorkspaceModel)

	for _, userWorkspace := range userWorkspaces {
		userIds = append(userIds, userWorkspace.UserID)
		userWorkspaceMap[userWorkspace.UserID] = userWorkspace
	}

	if len(userIds) == 0 {
		return []types.UserWorkspace{}, nil
	}

	users, err := prisma.Client.User.FindMany(
		db.User.ID.In(userIds),
	).Exec(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	var result []types.UserWorkspace
	for _, user := range users {
		userWorkspace, exists := userWorkspaceMap[user.ID]
		if exists {
			result = append(result, types.UserWorkspace{
				UserWorkspaceID: userWorkspace.ID,
				ID:              user.ID,
				Email:           user.Email,
				Name:            user.Name,
				ProfilePicture:  user.ProfilePicture,
				Role:            string(userWorkspace.Role), // Role from UserWorkspace
			})
		}
	}

	return result, nil
}

func (s *WorkspaceService) GetUserInWorkspace(userId, workspaceId string) (types.UserWorkspace, error) {
	// Find the user-workspace relation based on userId and workspaceId
	userWorkspace, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userId),
		db.UserWorkspace.WorkspaceID.Equals(workspaceId),
	).Exec(context.Background())

	if err != nil {
		return types.UserWorkspace{}, fmt.Errorf("failed to get user-workspace relation: %v", err)
	}

	if userWorkspace == nil {
		return types.UserWorkspace{}, fmt.Errorf("user not found in the specified workspace")
	}

	// Find the user based on userId
	user, err := prisma.Client.User.FindUnique(
		db.User.ID.Equals(userId),
	).Exec(context.Background())

	if err != nil {
		return types.UserWorkspace{}, fmt.Errorf("failed to get user: %v", err)
	}

	if user == nil {
		return types.UserWorkspace{}, fmt.Errorf("user not found")
	}

	// Return the user workspace details
	return types.UserWorkspace{
		UserWorkspaceID: userWorkspace.ID,
		ID:              user.ID,
		Email:           user.Email,
		Name:            user.Name,
		ProfilePicture:  user.ProfilePicture,
		Role:            string(userWorkspace.Role), // Assuming Role is stored in UserWorkspace
	}, nil
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

func (s *WorkspaceService) ChangeUserRole(workspaceId string, userId string, role db.UserRole) error {
	userwrk, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.WorkspaceID.Equals(workspaceId),
		db.UserWorkspace.UserID.Equals(userId),
	).Exec(context.Background())
	if err != nil {
		return err
	}
	_, err = prisma.Client.UserWorkspace.FindUnique(
		db.UserWorkspace.ID.Equals(userwrk.ID),
	).Update(
		db.UserWorkspace.Role.Set(role),
	).Exec(context.Background())
	if err != nil {
		panic(err)
	}
	return nil

}

func (s *WorkspaceService) KickUser(workspaceId string, userId string) (*db.WorkspaceModel, error) {
	userwrk, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userId),
		db.UserWorkspace.WorkspaceID.Equals(workspaceId),
	).With(db.UserWorkspace.Workspace.Fetch()).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	if userwrk.Workspace().OwnerID == userId {
		return nil, errors.New("Can't kick the owner of the workspace")
	}
	_, err = prisma.Client.UserWorkspace.FindUnique(
		db.UserWorkspace.ID.Equals(userwrk.ID),
	).Delete().Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return userwrk.Workspace(), nil
}

func (s *WorkspaceService) ChangeOwner(workspaceId, userId string) (*db.WorkspaceModel, error) {
	workspace, err := prisma.Client.Workspace.FindUnique(
		db.Workspace.ID.Equals(workspaceId),
	).Update(
		db.Workspace.Owner.Link(
			db.User.ID.Equals(userId),
		),
	).Exec(context.Background())

	return workspace, err
}

func (s *WorkspaceService) ChangeName(workspaceId, name string) (*db.WorkspaceModel, error) {
	workspace, err := prisma.Client.Workspace.FindUnique(
		db.Workspace.ID.Equals(workspaceId),
	).Update(
		db.Workspace.WorkspaceName.Set(name),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return workspace, err
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
