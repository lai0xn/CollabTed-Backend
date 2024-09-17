package router

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/config"
	"github.com/CollabTED/CollabTed-Backend/internal/sse"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func init() {
	// Initialize the middlware
	config.Load()
	echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWT_SECRET),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(types.Claims)
		},
	})
}

func SetRoutes(e *echo.Echo) {
	sse := sse.NewNotifier()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server Working check the docs at /swagger/index.html or the graphql playground at /graphql")
	})
	e.GET("/notifications", sse.NotificationHandler)
	v1 := e.Group("/api/v1")
	AuthRoutes(v1)
	OAuthRoutes(v1)
}
