package routes

import (
	"firebase.google.com/go/auth"
	"github.com/AniketSrivastava1/recruit/backend/internal/controllers"
	"github.com/AniketSrivastava1/recruit/backend/internal/middleware"
	"github.com/labstack/echo/v4"
)

func AddSubmissionRoutes(e *echo.Echo, authClient *auth.Client, ctrl *controllers.SubmissionController) {
	g := e.Group("/submission")
	g.Use(middleware.RequireFirebaseAuth(authClient))

	g.POST("/submit", ctrl.Submit)
	g.GET("/:id/status", ctrl.Status)
	g.GET("/:id/details", ctrl.Details)
	g.GET("/list", ctrl.List)
}
