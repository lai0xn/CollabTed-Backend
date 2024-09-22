package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

type calendarHandler struct {
	srv services.EventService
}

func NewCalendarHandler() *calendarHandler {
	return &calendarHandler{
		srv: *services.NewEventService(),
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

	// Call service to create event
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

	// Call service to list events by workspace
	data, err := h.srv.ListEventsByWorkspace(workspaceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, data)
}
