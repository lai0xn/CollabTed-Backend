package middlewares

import (
	"fmt"
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/config"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		logger.LogDebug().Msg(fmt.Sprintf("Cookie value: %s", cookie.Value))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid cookie")
		}

		tokenString := cookie.Value
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
		}

		token, err := jwt.ParseWithClaims(tokenString, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWT_SECRET), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		claims, ok := token.Claims.(*types.Claims)
		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		c.Set("user", claims)
		logger.LogDebug().Msg(fmt.Sprintf("User: %+v", claims))
		return next(c)
	}
}
