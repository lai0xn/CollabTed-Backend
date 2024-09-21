package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

type workspaceHandler struct {
	srv services.WorkspaceService
}

func NewWorkspaceHandler() *workspaceHandler {
	return &workspaceHandler{
		srv: *services.NewWorkspaceService(),
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
	data, err := h.srv.CreateWorkspace(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, data)
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
func (h *workspaceHandler) GetWorkspace(c echo.Context) error {
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

// CreateInvitation example
//
//	@Summary	Create an invitation
//	@Tags		workspace
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Param		body		body		types.Invitation	true	"Invitation details"
//	@Success	200		{object}	types.Invitation
//	@Security	BearerAuth
//	@Router		/invitations [post]
func (h *workspaceHandler) CreateInvitation(c echo.Context) error {
	var payload types.InviteUserRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	invitation, err := services.CreateInvitation(payload.Email, payload.WorkspaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, invitation)
}

func (h *workspaceHandler) AcceptInvitation(c echo.Context) error {
	token := c.QueryParam("token")

	invitation, err := services.AcceptInvitation(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, invitation)
}

func (h *workspaceHandler) DeclineInvitation(c echo.Context) error {
	token := c.QueryParam("token")

	invitation, err := services.DeclineInvitation(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, invitation)
}
