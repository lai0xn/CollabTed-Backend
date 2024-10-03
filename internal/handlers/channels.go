package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

type channelHandler struct {
	srv services.ChannelService
}

func NewChannelHandler() *channelHandler {
	return &channelHandler{
		srv: *services.NewChannelService(),
	}
}

func (h *channelHandler) GetWorkspaceChannels(c echo.Context) error {
	worksapceID := c.Param("id")
	channels, err := h.srv.ListChannelsByWorkspace(worksapceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, channels)
}

func (h *channelHandler) GetChannel(c echo.Context) error {
	worksapceID := c.Param("id")
	channel, err := h.srv.GetChannelById(worksapceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, channel)
}

func (h *channelHandler) CreateChannel(c echo.Context) error {
	var data types.ChannelD
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	globalRoomJoinToken, err := h.srv.CreateChannel(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, globalRoomJoinToken)
}

func (h *channelHandler) AddParticipant(c echo.Context) error {
	var data types.ParticipantD
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	participant, err := h.srv.AddParticipant(data.WorkspaceID, data.ChannelID, data.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, participant)
}
