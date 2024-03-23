package auth

import (
	token "gandalf-data-aggregator/pkg/jwt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtMaker token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing auth token")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth format")
			}

			tokenString := parts[1]
			payload, err := jwtMaker.VerifyToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			if payload.ExpiredAt.Before(time.Now()) {
				return echo.NewHTTPError(http.StatusUnauthorized, "token has expired")
			}

			c.Set("UserID", payload.UserID)
			return next(c)
		}
	}
}
