package middleware

import (
	"net/http"

	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
	"github.com/labstack/echo/v4"
)

func RequireAdminRole(adminStore *stores.AdminStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if dummyAuthEnabled() {
				return next(c)
			}
			userId := c.Get("userId").(string)
			isAdmin, err := adminStore.IsAdmin(c.Request().Context(), userId)
			if err != nil || !isAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "Admin privileges required")
			}
			return next(c)
		}
	}
}
