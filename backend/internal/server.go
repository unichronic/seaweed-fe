package internal

import (
	"context"

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

	return e
}

func StartEchoServer(lifecycle fx.Lifecycle, e *echo.Echo, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Echo server on :8080")
			go e.Start(":8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Echo server")
			return e.Shutdown(ctx)
		},
	})
}
