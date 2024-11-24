package router

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/config"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/CollabTED/CollabTed-Backend/internal/sse"
	"github.com/CollabTED/CollabTed-Backend/internal/ws"
	"github.com/labstack/echo/v4"
)

func init() {
	// Initialize the middlware
	config.Load()
}

func SetRoutes(e *echo.Echo) {
	sse := sse.NewNotifier()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server Working check the docs at /swagger/index.html or the graphql playground at /graphql")
	})

	e.GET("/ws", ws.WsChatHandler{}.Chat, middlewares.AuthMiddleware)
	e.GET("/notifications", sse.NotificationHandler)

	v1 := e.Group("/api/v1")
	AuthRoutes(v1)
	OAuthRoutes(v1)
	WorkspaceRoutes(v1)
	CalendarRoutes(v1)
	CallsRoutes(v1)
	ChannelsRoutes(v1)
	MessageRoutes(v1)
	ProjectsRoutes(v1)
	StatusRoutes(v1)
	BoardRoutes(v1)
	TasksRoutes(v1)
	LiveBoardRoutes(v1)
}
