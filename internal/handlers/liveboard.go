package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

type liveBoardHandler struct {
	srv services.LiveBoardService
}

func NewLiveBoardHandler() *liveBoardHandler {
	return &liveBoardHandler{
		srv: *services.NewLiveBoardService(),
	}
}

func (h *liveBoardHandler) GetBoard(c echo.Context) error {
	boardId := c.Param("boardId")
	data, err := h.srv.GetBoard(boardId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *liveBoardHandler) GetWorkspaceBoards(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	data, err := h.srv.GetWorkspaceBoards(workspaceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}



func (h *liveBoardHandler) CreateBoard(c echo.Context) error {
	var data types.LiveBoardD
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	message, err := h.srv.CreateBoard(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, message)
}

func (h *liveBoardHandler) DeleteBoard(c echo.Context) error {
	boardId := c.Param("boardId")
	result,err := h.srv.DeleteBoard(boardId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
