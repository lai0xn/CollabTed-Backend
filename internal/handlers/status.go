package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
	"github.com/labstack/echo/v4"
)

type StatusHandler struct {
	statusService  services.StatusService
	projectService services.ProjectService
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{
		statusService:  *services.NewStatusService(),
		projectService: *services.NewProjectService(),
	}
}

func (h *StatusHandler) CreateStatus(c echo.Context) error {
	var statusD types.StatusD
	err := c.Bind(&statusD)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	status, err := h.statusService.CreateStatus(statusD)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusCreated, status)
}

func (h *StatusHandler) GetStatusesByProject(c echo.Context) error {
	projectID := c.Param("projectId")
	claims := c.Get("user").(*types.Claims)

	statuses, err := h.statusService.GetStatusesByProject(projectID, claims.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, statuses)
}

func (h *StatusHandler) GetStatusByID(c echo.Context) error {
	statusID := c.Param("statusId")
	claims := c.Get("user").(*types.Claims)
	status, err := h.statusService.GetStatusByID(statusID, claims.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

func (h *StatusHandler) DeleteStatus(c echo.Context) error {
	statusID := c.Param("statusId")
	WorksapceID := c.Param("workspaceId")
	claims := c.Get("user").(*types.Claims)

	canCreate, err := h.projectService.CanUserPerformAction(claims.ID, WorksapceID, db.UserRoleAdmin)
	if err != nil || !canCreate {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create projects")
	}

	err = h.statusService.DeleteStatus(statusID, claims.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *StatusHandler) EditStatus(c echo.Context) error {
	statusID := c.Param("statusId")
	WorksapceID := c.Param("workspaceId")
	claims := c.Get("user").(*types.Claims)

	canCreate, err := h.projectService.CanUserPerformAction(claims.ID, WorksapceID, db.UserRoleAdmin)
	if err != nil || !canCreate {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create projects")
	}

	var statusD types.StatusD
	err = c.Bind(&statusD)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := h.statusService.EditStatus(statusID, claims.ID, statusD)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
