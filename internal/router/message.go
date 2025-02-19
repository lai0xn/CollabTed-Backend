package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func MessageRoutes(e *echo.Group) {
	h := handlers.NewMessageHandler()
	messages := e.Group("/messages", middlewares.AuthMiddleware)
	messages.POST("/", h.SendMessage)
	messages.POST("/pin/:messageId", h.PinMessage)
	messages.GET("/:channelId/pinned", h.GetPinnedMessages)
	messages.GET("/:channelId", h.GetMessages)
	messages.GET("/attachments/:channelId", h.GetAttachments)
	messages.DELETE("/:messageId", h.DeleteMessage)
	messages.DELETE("/:attachmentId", h.DeleteMessage)
	messages.POST("/attachment", h.UploadAttachment)
	messages.DELETE("/attachment/:id", h.DeleteAttachment)
}
