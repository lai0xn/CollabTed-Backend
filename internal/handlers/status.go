package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

type StatusHandler struct {
	statusService services.StatusService
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{*services.NewStatusService()}
}

func (h *StatusHandler) CreateStatus(c echo.Context) error {
	var statusD types.StatusD
	err := c.Bind(&statusD)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	status, err := h.statusService.CreateStatus(statusD, c.Get("userID").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusCreated, status)
}

func (h *StatusHandler) GetStatusesByProject(c echo.Context) error {
	projectID := c.Param("projectID")

	statuses, err := h.statusService.GetStatusesByProject(projectID, c.Get("userID").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, statuses)
}

func (h *StatusHandler) GetStatusByID(c echo.Context) error {
	statusID := c.Param("statusID")

	status, err := h.statusService.GetStatusByID(statusID, c.Get("userID").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

func (h *StatusHandler) DeleteStatus(c echo.Context) error {
	statusID := c.Param("statusID")

	err := h.statusService.DeleteStatus(statusID, c.Get("userID").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
