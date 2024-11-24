package server

import (
	"github.com/CollabTED/CollabTed-Backend/config"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/CollabTED/CollabTed-Backend/internal/router"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	PORT string
}

func NewServer(port string) *Server {
	return &Server{
		PORT: port,
	}
}

func (s *Server) Setup(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Static("/public", "public")
	router.SetRoutes(e)

	// CORS configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{config.ALLOWED_ORIGINS},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH ,echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// logging middleware
	e.Use(middlewares.LoggingMiddleware)
}

func (s *Server) Run() {
	e := echo.New()
	s.Setup(e)
	logger.LogInfo().Msg("graphql server running on /graphql")
	logger.LogInfo().Msg(e.Start(s.PORT).Error())
}
