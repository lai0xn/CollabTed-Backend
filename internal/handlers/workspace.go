package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
	"github.com/labstack/echo/v4"
)

type workspaceHandler struct {
	srv      services.WorkspaceService
	boardSrv services.BoardService
}

func NewWorkspaceHandler() *workspaceHandler {
	return &workspaceHandler{
		srv:      *services.NewWorkspaceService(),
		boardSrv: *services.NewBoardService(),
	}
}

// CreateWorkspace example
//
//	@Summary	Create a new workspace
//	@Tags		workspace
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Param		body		body		types.WorkspaceD	true	"Workspace details"
//	@Success	201		{object}	types.WorkspaceD
//	@Router		/workspaces [post]
func (h *workspaceHandler) CreateWorkspace(c echo.Context) error {
	var payload types.WorkspaceD
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	workspace, err := h.srv.CreateWorkspace(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err = h.boardSrv.SaveBoard(types.BoardD{
		WorkspaceID: workspace.ID,
		Elements:    []json.RawMessage{},
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := c.Get("user").(*types.Claims)

	_, err = prisma.Client.UserWorkspace.CreateOne(
		db.UserWorkspace.User.Link(
			db.User.ID.Equals(claims.ID),
		),
		db.UserWorkspace.Workspace.Link(
			db.Workspace.ID.Equals(workspace.ID),
		),
		db.UserWorkspace.Role.Set(db.UserRoleAdmin),
		db.UserWorkspace.JoinedAt.Set(time.Now()),
	).Exec(context.Background())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to associate user with workspace: "+err.Error())
	}

	return c.JSON(http.StatusCreated, workspace)
}

// GetWorkspaces example
//
//	@Summary	List workspaces for the authenticated user
//	@Tags		workspace
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Success	200		{array}		types.WorkspaceD
//	@Security	BearerAuth
//	@Router		/workspaces [get]
func (h *workspaceHandler) GetWorkspaces(c echo.Context) error {
	if c.Get("user") == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	claims := c.Get("user").(*types.Claims)

	data, err := h.srv.ListWorkspaces(claims.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

// Get workspace by id
//
//	@Summary	Get workspace by id
//	@Tags		workspace
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Workspace id"
//	@Success	200	{object}	types.WorkspaceD
//	@Router		/workspaces/{id} [get]
func (h *workspaceHandler) GetWorkspaceById(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	data, err := h.srv.GetWorkspaceById(workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

// InviteUser example
//
//	@Summary	Invite a user to a workspace
//	@Tags		workspace
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Param		body		body		types.InviteUserD	true	"Invitation details"
//	@Success	200		{string}	string	"Invitation sent successfully"
//	@Router		/workspaces/invite [post]
func (h *workspaceHandler) InviteUser(c echo.Context) error {
	var payload types.InviteUserD
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := c.Get("user").(*types.Claims)
	canInvite, err := h.srv.CanUserPerformAction(claims.ID, payload.WorkspaceID, db.UserRoleAdmin)
	if err != nil || !canInvite {
		canInvite, err = h.srv.CanUserPerformAction(claims.ID, payload.WorkspaceID, db.UserRoleManager)
	}
	if err != nil || !canInvite {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to invite users to this workspace")
	}

	err = h.srv.SendInvitation(payload.Email, payload.WorkspaceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send invitation: "+err.Error())
	}

	return c.JSON(http.StatusOK, "Invitation sent successfully")
}

func (h *workspaceHandler) AcceptInvitation(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Token is required")
	}

	claims := c.Get("user").(*types.Claims)
	err := h.srv.AcceptInvitation(claims.ID, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Successfully joined the workspace")
}

func (h *workspaceHandler) GetAllUsersInWorkspace(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	data, err := h.srv.GetAllUsersInWorkspace(workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *workspaceHandler) GetUserInWorkspace(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	userId := c.Param("userId")
	data, err := h.srv.GetUserInWorkspace(workspaceId, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *workspaceHandler) GetAllInvites(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	invitations, err := h.srv.GetInvitations(workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, invitations)
}

func (h *workspaceHandler) DeleteInvitation(c echo.Context) error {
	invitationId := c.Param("invitationId")
	err := h.srv.DeleteInvitation(invitationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Invitation deleted successfully")
}
