
package router

import (
	"github.com/CollabTED/CollabTed-Backend/internal/handlers"
	middlewares "github.com/CollabTED/CollabTed-Backend/internal/middlewares/rest"
	"github.com/labstack/echo/v4"
)

func LiveBoardRoutes(e *echo.Group) {
	h := handlers.NewLiveBoardHandler()

	board := e.Group("/liveboard", middlewares.AuthMiddleware)
	board.GET("/:boardId", h.GetBoard)
	board.GET("/workspace:/workspaceId",h.GetWorkspaceBoards)
	board.POST("/", h.CreateBoard)
	board.DELETE("/:boardId", h.DeleteBoard)
}
