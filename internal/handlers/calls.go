package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/labstack/echo/v4"
)

type callHandler struct {
	srv services.CallService
}

func NewCallHandler() *callHandler {
	return &callHandler{
		srv: *services.NewCallService(),
	}
}

func (h *callHandler) GetGlobalJoinToken(c echo.Context) error {
	participantName := c.Param("participantName")
	workspaceId := c.Param("workspaceId")

	globalRoomJoinToken, err := h.srv.GetGlobalJoinToken(participantName, workspaceId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, globalRoomJoinToken)
}
