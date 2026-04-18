package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

func RequireFirebaseAuth(authClient *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
			}

			idToken := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
			token, err := authClient.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid ID Token")
			}

			c.Set("user", token)
			c.Set("userId", token.UID)
			return next(c)
		}
	}
}

func OptionalFirebaseAuth(authClient *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c)
			}

			idToken := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
			token, err := authClient.VerifyIDToken(context.Background(), idToken)
			if err == nil {
				c.Set("user", token)
				c.Set("userId", token.UID)
			}
			return next(c)
		}
	}
}
