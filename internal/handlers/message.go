package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

type messageHandler struct {
	srv services.MessageService
}

func NewMessageHandler() *messageHandler {
	return &messageHandler{
		srv: *services.NewMessageService(),
	}
}

func (h *messageHandler) GetMessages(c echo.Context) error {
	channelId := c.Param("channelId")
	data, err := h.srv.GetMessagesByChannel(channelId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *messageHandler) SendMessage(c echo.Context) error {
	var data types.MessageD
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	message, err := h.srv.SendMessage(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, message)
}

func (h *messageHandler) DeleteMessage(c echo.Context) error {
	messageId := c.Param("messageId")
	err := h.srv.DeleteMessage(messageId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Message deleted successfully")
}
