package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

type BoardHandler struct {
	srv services.BoardService
}

func NewBoardHandler() *BoardHandler {
	return &BoardHandler{
		srv: *services.NewBoardService(),
	}
}


func (h *BoardHandler) UpdateBoard(c echo.Context) error {
	var boardId = c.Param("boardId")

	var request types.BoardD
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := h.srv.UpdateBoard(request, boardId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (h *BoardHandler) GetBoard(c echo.Context) error {
	workspaceId := c.Param("workspaceId")

	elements, err := h.srv.GetBoard(workspaceId)

	if err != nil {
		return c.JSON(http.StatusNotFound, "Board not found")
	}

	return c.JSON(http.StatusOK, elements)
}
func (h *BoardHandler) CreateBoard(c echo.Context) error {
	var request types.BoardD
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := h.srv.SaveBoard(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to save board")
	}

	return c.JSON(http.StatusOK, result)
}
