package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
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

	globalRoomJoinToken, roomId, err := h.srv.GetGlobalJoinToken(participantName, workspaceId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := map[string]interface{}{
		"roomId": roomId,
		"token":  globalRoomJoinToken,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *callHandler) GetPrivatelJoinToken(c echo.Context) error {
	Caller := c.Param("participantName")
	Receiver := c.Param("receiverId")
	workspaceId := c.Param("workspaceId")

	privateRoomJoinToken, roomId, err := h.srv.GetPrivateJoinToken(Caller, workspaceId)

	response := map[string]interface{}{
		"roomId": roomId,
		"token":  privateRoomJoinToken,
	}

	logger.LogInfo().Msg(Receiver)
	logger.LogInfo().Msg(privateRoomJoinToken)
	// TODO: add the logic to send the private room join token to the requester

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *callHandler) JoinRoomToken(c echo.Context) error {
	roomId := c.Param("roomId")
	participantName := c.Param("participantName")

	joinRoomToken, err := h.srv.JoinRoomToken(roomId, participantName)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, joinRoomToken)
}
