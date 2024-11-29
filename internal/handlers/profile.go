package handlers

import (
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	srv *services.ProfileService
}

func NewProfileHandler()*ProfileHandler{
	return &ProfileHandler{
		srv: services.NewProfileService(),
	}
}

func (h *ProfileHandler)UpdateProfile(c echo.Context) error{
	var paylaod types.ProfileUpdate
	usr := c.Get("user").(*types.Claims)
	if err := c.Bind(&paylaod);err!= nil {
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}

	user,err := h.srv.UpdateUser(usr.ID,paylaod)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	return c.JSON(http.StatusOK,user)
}

func (h *ProfileHandler) GetProfile(c echo.Context) error {return nil}
