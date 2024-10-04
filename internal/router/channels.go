package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func ChannelsRoutes(e *echo.Group) {
	h := handlers.NewChannelHandler()

	channels := e.Group("/channels", middlewares.AuthMiddleware)
	channels.GET("/:workspaceId", h.GetChannel)
	channels.POST("/", h.CreateChannel)
	channels.GET("/worksapce/:workspaceId", h.GetWorkspaceChannels)
	channels.POST("/participants/add", h.AddParticipant)
}
