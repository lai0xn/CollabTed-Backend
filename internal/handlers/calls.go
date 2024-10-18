package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/internal/sse"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/labstack/echo/v4"
)

type callHandler struct {
	srv      services.CallService
	notifier *sse.Notifier
}

func NewCallHandler() *callHandler {
	return &callHandler{
		srv:      *services.NewCallService(),
		notifier: sse.NewNotifier(),
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
	CallerId := c.QueryParam("callerId")
	ReceiverId := c.Param("receiverId")
	workspaceId := c.Param("workspaceId")

	privateRoomJoinToken, roomId, err := h.srv.GetPrivateJoinToken(Caller, workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := map[string]interface{}{
		"roomId": roomId,
		"token":  privateRoomJoinToken,
	}

	logger.LogInfo().Msg(ReceiverId)
	logger.LogInfo().Msg(privateRoomJoinToken)

	err = h.notifier.NotifyUser(ReceiverId, roomId, CallerId)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to notify user"})
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
