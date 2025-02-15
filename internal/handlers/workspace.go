package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/internal/sse"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
	"github.com/labstack/echo/v4"
)

type workspaceHandler struct {
	srv      services.WorkspaceService
	csrv     services.ChannelService
	boardSrv services.BoardService
	notifier *sse.Notifier
}

func NewWorkspaceHandler() *workspaceHandler {
	return &workspaceHandler{
		srv:      *services.NewWorkspaceService(),
		csrv:     *services.NewChannelService(),
		boardSrv: *services.NewBoardService(),
		notifier: sse.NewNotifier(),
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
		canInvite, err = h.srv.CanUserPerformAction(claims.ID, payload.WorkspaceID, db.UserRoleAdmin)
	}
	if err != nil || !canInvite {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to invite users to this workspace")
	}

	err = h.srv.SendInvitation(payload.Email, payload.WorkspaceID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Failed to send invitation: "+err.Error())
	}

	return c.JSON(http.StatusOK, "Invitation sent successfully")
}

func (h *workspaceHandler) AcceptInvitation(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Token is required")
	}

	claims := c.Get("user").(*types.Claims)
	workspaceID, err := h.srv.AcceptInvitation(claims.ID, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	logger.LogDebug().Msgf("?????Workspace ID: %s", workspaceID)

	err = h.notifier.NotifyJoinUser(claims.ID, workspaceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	logger.LogDebug().Msgf("??????Invitation accepted for user ID: %s", claims.ID)

	return c.JSON(http.StatusOK, "Successfully joined the workspace")
}

func (h *workspaceHandler) GetWorkspaceChannels(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	data, err := h.csrv.ListChannelsByWorkspace(workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *workspaceHandler) GetAllUsersInWorkspace(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	data, err := h.srv.GetAllUsersInWorkspace(workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *workspaceHandler) GetConnectedUsers(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	data, err := h.srv.GetConnectedUsers(workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *workspaceHandler) GetUserInWorkspace(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	userId := c.Param("userId")
	data, err := h.srv.GetUserInWorkspace(userId, workspaceId)
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

func (h *workspaceHandler) ChangeName(c echo.Context) error {
	type payload struct {
		Name string `json:"name"`
	}
	claims := c.Get("user").(*types.Claims)

	worksapceId := c.Param("workspaceId")
	var data payload
	canPerform, err := h.srv.CanUserPerformAction(claims.ID, worksapceId, db.UserRoleAdmin)

	if err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if !canPerform {
		return echo.NewHTTPError(http.StatusForbidden, "you are not authorized to perform this action")
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	workspace, err := h.srv.ChangeName(worksapceId, data.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, workspace)
}

func (h *workspaceHandler) ChangeOwner(c echo.Context) error {
	type payload struct {
		UserID string `json:"userId"`
	}
	claims := c.Get("user").(*types.Claims)

	worksapceId := c.Param("workspaceId")
	var data payload
	canPerform, err := h.srv.CanUserPerformAction(claims.ID, worksapceId, db.UserRoleAdmin)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if !canPerform {
		return echo.NewHTTPError(http.StatusForbidden, "you are not authorized to perform this action")
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	workspace, err := h.srv.ChangeOwner(worksapceId, claims.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, workspace)

}

func (h *workspaceHandler) KickUser(c echo.Context) error {
	userId := c.Param("userId")
	workspaceID := c.Param("workspaceId")

	logger.LogInfo().Msg("TRYING")
	err := h.notifier.NotifyKickUser(userId, workspaceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	logger.LogInfo().Msg("OK")

	workspace, err := h.srv.KickUser(workspaceID, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, workspace)
}

func (h *workspaceHandler) ChangeUserRole(c echo.Context) error {
	userId := c.Param("userId")
	workspaceID := c.Param("workspaceId")
	type payload struct {
		Role string `json:"role"`
	}
	var role payload
	if err := c.Bind(&role); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	claims := c.Get("user").(*types.Claims)
	canPerform, err := h.srv.CanUserPerformAction(claims.ID, workspaceID, db.UserRoleAdmin)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !canPerform {
		return echo.NewHTTPError(http.StatusForbidden, "you are not authorized to perform this action")
	}

	if err := h.srv.ChangeUserRole(workspaceID, userId, db.UserRole(role.Role)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "you don't have permission to perform this action",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"msg": "role changed",
	})
}
