package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/mail"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	srv      *services.AuthService
	verifier *mail.EmailVerifier
}

func NewAuthHandler() *authHandler {
	return &authHandler{
		srv:      services.NewAuthService(),
		verifier: mail.NewVerifier(),
	}
}

// Login example
//
//	@Summary	Login endpoint
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body		body		types.LoginPayload	true	"Login details"
//	@Router		/auth/login [post]
func (h *authHandler) Login(c echo.Context) error {
	var payload types.LoginPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := validate.Struct(payload)
	if err != nil {
		e := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, utils.NewValidationError(e))
	}

	user, err := h.srv.CheckUser(payload.Email, payload.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := utils.SetJWTAsCookie(c.Response().Writer, user.ID, user.Email, user.Name, user.ProfilePicture); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to set JWT cookie")
	}

	return c.JSON(http.StatusOK, types.Response{
		"message": "token set",
	})
}

// Registration example
//
//	@Summary	Registration endpoint
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body		body		types.RegisterPayload	true	"Registration details"
//	@Router		/auth/register [post]
func (h *authHandler) Register(c echo.Context) error {
	var payload types.RegisterPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := validate.Struct(payload)
	if err != nil {
		e := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, utils.NewValidationError(e))
	}

	if payload.ProfilePicture == "" {
		avatarURL := fmt.Sprintf(
			strings.ReplaceAll("https://ui-avatars.com/api/?name=%s", " ", ""),
			url.QueryEscape(payload.Name),
		)

		payload.ProfilePicture = avatarURL
	}

	user, err := h.srv.CreateUser(payload.Name, payload.Email, payload.Password, payload.ProfilePicture)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.verifier.SendVerfication(user.ID, []string{user.Email})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, types.Response{
		"message": "verification email sent",
		"userID":  user.ID,
	})
}

// Email verification example
//
//	@Summary	Verification endpoint
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		id		query		string	true	"userid"
//	@Param		otp		query		string	true	"otp"
//	@Router		/auth/verify [post]
func (h *authHandler) VerifyUser(c echo.Context) error {
	id := c.QueryParam("id")
	otp := c.QueryParam("otp")
	if err := h.verifier.Verify(id, otp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := h.srv.ActivateUser(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err := h.srv.ActivateUser(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, types.Response{
		"message": "user activated",
	})
}

func (h *authHandler) CheckUser(c echo.Context) error {
	if c.Get("user") == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "authenticated"})
}

func (h *authHandler) Me(c echo.Context) error {
	if c.Get("user") == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	claims := c.Get("user").(*types.Claims)

	utils.FetchAndEncodeImageToBase64(claims.ProfilePicture)

	return c.JSON(http.StatusOK, claims)
}

func (h *authHandler) Logout(c echo.Context) error {
	if err := utils.DeleteJWTCookie(c.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to delete JWT cookie")
	}
	return c.JSON(http.StatusOK, types.Response{
		"message": "token deleted",
	})
}

func (h *authHandler) SendRessetLink(c echo.Context) error {
	email := c.QueryParam("email")
	if err := h.srv.SendRessetLink(email); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "resset link sent",
	})
}

func (h *authHandler) RessetPassword(c echo.Context) error {
	var body struct {
		Password string
		Email    string
		Token    string
	}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.srv.RessetPassword(body.Email, body.Token, body.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "resset link sent",
	})
}
