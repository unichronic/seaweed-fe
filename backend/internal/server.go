package internal

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewEchoServer(logger *zap.Logger) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return e
}

func StartEchoServer(lifecycle fx.Lifecycle, e *echo.Echo, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}
			logger.Info("Starting Echo server", zap.String("port", port))
			go func() {
				if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
					logger.Error("Echo server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Echo server")
			return e.Shutdown(ctx)
		},
	})
}
