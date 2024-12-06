package middlewares

import (
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type MiddlewareFunc func(next echo.HandlerFunc) echo.HandlerFunc

type Protected interface {
	AccessPerm(userId string) error
	WritePerm(userId string) error
}

func MiddlewareBuidler(p Protected) (echo.MiddlewareFunc, echo.MiddlewareFunc) {
	readM := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(types.Claims)
			if err := p.AccessPerm(claims.ID); err != nil {
				return echo.NewHTTPError(echo.ErrForbidden.Code)

			}
			return next(c)
		}
	}
	writeM := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(types.Claims)
			if err := p.WritePerm(claims.ID); err != nil {
				return echo.NewHTTPError(echo.ErrForbidden.Code)

			}
			return next(c)
		}
	}
	return readM, writeM
}
