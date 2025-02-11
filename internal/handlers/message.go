package handlers

import (
	"net/http"
	"strconv"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/cloudinary"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
	page := c.QueryParam("p")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	data, err := h.srv.GetMessagesByChannel(channelId, pageInt)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *messageHandler) GetPinnedMessages(c echo.Context) error {
	channelId := c.Param("channelId")
	page := c.QueryParam("p")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	data, err := h.srv.GetPinnedMessages(channelId, pageInt)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *messageHandler) GetAttachments(c echo.Context) error {
	channelId := c.Param("channelId")
	data, err := h.srv.GetAttachmentsByChannel(channelId)
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
	messageId := c.Param("messageID")
	err := h.srv.DeleteMessage(messageId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Message deleted successfully")
}
func (h *messageHandler) DeleteAttachment(c echo.Context) error {
	attachmentID := c.Param("attachmentID")
	err := h.srv.DeleteMessage(attachmentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Message deleted successfully")
}

func (s *messageHandler) UploadAttachment(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File is required")
	}
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to open file")
	}
	defer src.Close()

	channelID := c.FormValue("channelID")
	if channelID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "channelID is required")
	}
	senderID := c.FormValue("senderID")
	if senderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "senderID is required")
	}
	workspaceID := c.FormValue("workspaceID")
	if workspaceID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "workspaceID is required")
	}

	uploadParams := uploader.UploadParams{Folder: "attachments"}
	result, err := cloudinary.GetUploader().Upload(c.Request().Context(), src, uploadParams)
	if err != nil {
		logger.Logger.Err(err).Msg("Failed to upload file to Cloudinary")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload file")
	}

	attachment := types.AttachmentD{
		ChannelID:   channelID,
		SenderID:    senderID,
		WorkspaceID: workspaceID,
		File:        result.SecureURL,
		Title:       file.Filename,
	}

	_, err = s.srv.CreateAttachment(attachment)
	if err != nil {
		logger.Logger.Err(err).Msg("Failed to save attachment in database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload file")
	}
	
	return c.JSON(http.StatusOK, attachment)
}

func (h *messageHandler) PinMessage(c echo.Context) error {
	messageId := c.Param("messageId")
	err := h.srv.PingMessage(messageId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "message pinned",
	})
}
