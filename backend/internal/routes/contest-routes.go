package routes

import (
	"firebase.google.com/go/auth"
	"github.com/AniketSrivastava1/recruit/backend/internal/controllers"
	"github.com/AniketSrivastava1/recruit/backend/internal/middleware"
	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
	"github.com/labstack/echo/v4"
)

func AddContestRoutes(e *echo.Echo, authClient *auth.Client, ctrl *controllers.ContestController) {
	e.GET("/contests/list", ctrl.ListContests)
}

func AddAdminRoutes(e *echo.Echo, authClient *auth.Client, adminStore *stores.AdminStore, ctrl *controllers.ContestController) {
	g := e.Group("/admin")
	g.Use(middleware.RequireFirebaseAuth(authClient))
	g.Use(middleware.RequireAdminRole(adminStore))

	g.POST("/contest", ctrl.CreateContest)
	g.POST("/:contestId/problem", ctrl.AddProblem)
}
