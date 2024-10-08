package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func MessageRoutes(e *echo.Group) {
	h := handlers.NewMessageHandler()
	messages := e.Group("/messages", middlewares.AuthMiddleware)

	messages.POST("/send", h.SendMessage)
	messages.GET("/:channelId/get", h.GetMessages)
	messages.POST("/:messageId/delete", h.DeleteMessage)
}
