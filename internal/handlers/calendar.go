package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
	"github.com/labstack/echo/v4"
)

type calendarHandler struct {
	srv          services.EventService
	workspaceSrv services.WorkspaceService
}

func NewCalendarHandler() *calendarHandler {
	return &calendarHandler{
		srv:          *services.NewEventService(),
		workspaceSrv: *services.NewWorkspaceService(),
	}
}

// CreateEvent example
//
//	@Summary	Create a new event
//	@Tags		event
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Param		body		body		types.EventD	true	"Event details"
//	@Success	201		{object}	types.EventD
//	@Router		/events [post]
func (h *calendarHandler) CreateEvent(c echo.Context) error {
	var payload types.EventD
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := c.Get("user").(*types.Claims)

	canCreateAdmin, err := h.workspaceSrv.CanUserPerformAction(claims.ID, payload.WorkspaceID, db.UserRoleAdmin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error checking user permissions: "+err.Error())
	}

	canCreateManager, err := h.workspaceSrv.CanUserPerformAction(claims.ID, payload.WorkspaceID, db.UserRoleManager)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error checking user permissions: "+err.Error())
	}

	if !canCreateAdmin && !canCreateManager {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create an event in this workspace")
	}

	data, err := h.srv.CreateEvent(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, data)
}

// ListEvents example
//
//	@Summary	List events for a workspace
//	@Tags		event
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Param		workspaceId		path		string	true	"Workspace ID"
//	@Success	200		{array}		types.EventD
//	@Router		/workspaces/{workspaceId}/events [get]
func (h *calendarHandler) ListEvents(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "workspaceId is required")
	}

	data, err := h.srv.ListEventsByWorkspace(workspaceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, data)
}
