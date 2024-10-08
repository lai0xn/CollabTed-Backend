package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
	"github.com/labstack/echo/v4"
)

type projectHandler struct {
	srv services.ProjectService
}

func NewProjectHandler() *projectHandler {
	return &projectHandler{
		srv: *services.NewProjectService(),
	}
}

// CreateProject example
//
//	@Summary	Create a new project
//	@Tags		project
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Param		body		body		types.ProjectD	true	"Project details"
//	@Success	201		{object}	types.ProjectD
//	@Router		/projects [post]
func (h *projectHandler) CreateProject(c echo.Context) error {
	var payload types.ProjectD
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := c.Get("user").(*types.Claims)

	// Check if the user can create projects
	canCreate, err := h.srv.CanUserPerformAction(claims.ID, payload.WorksapceID, db.UserRoleManager)
	if err != nil || !canCreate {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create projects")
	}

	project, err := h.srv.CreateProject(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, project)
}

// GetProjects example
//
//	@Summary	List projects for the authenticated user
//	@Tags		project
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Success	200		{array}		types.ProjectD
//	@Router		/projects [get]
func (h *projectHandler) GetProjects(c echo.Context) error {
	claims := c.Get("user").(*types.Claims)

	data, err := h.srv.ListProjectsByWorkspace(c.Param("worksapceID"), claims.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

// GetProjectById example
//
//	@Summary	Get project by id
//	@Tags		project
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Project id"
//	@Success	200	{object}	types.ProjectD
//	@Router		/projects/{id} [get]
func (h *projectHandler) GetProjectById(c echo.Context) error {
	projectId := c.Param("id")
	data, err := h.srv.GetProjectById(projectId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}
