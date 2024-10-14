package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func BoardRoutes(e *echo.Group) {
	h := handlers.NewBoardHandler()
	boards := e.Group("/boards", middlewares.AuthMiddleware)
	boards.PUT("/update/:boardId", h.UpdateBoard)
	boards.GET("/list/:workspaceId", h.GetBoard)
}
