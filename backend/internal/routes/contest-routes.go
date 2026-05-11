package routes

import (
	"firebase.google.com/go/auth"
	"github.com/AniketSrivastava1/recruit/backend/internal/controllers"
	"github.com/AniketSrivastava1/recruit/backend/internal/middleware"
	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
	"github.com/labstack/echo/v4"
)

func AddContestRoutes(e *echo.Echo, authClient *auth.Client, ctrl *controllers.ContestController) {
	e.GET("/contests/list", ctrl.ListContests, middleware.OptionalFirebaseAuth(authClient))
	e.GET("/contests/:id", ctrl.GetContest, middleware.OptionalFirebaseAuth(authClient))
	e.POST("/contests/:id/registration", ctrl.Register, middleware.RequireFirebaseAuth(authClient))
	e.GET("/contests/:id/problems", ctrl.ListProblems, middleware.RequireFirebaseAuth(authClient))
	e.GET("/contests/:id/problems/:problemId", ctrl.GetProblem, middleware.RequireFirebaseAuth(authClient))
	e.GET("/contests/:id/leaderboard", ctrl.GetLeaderboard, middleware.OptionalFirebaseAuth(authClient))
}

func AddAdminRoutes(e *echo.Echo, authClient *auth.Client, adminStore *stores.AdminStore, ctrl *controllers.ContestController) {
	g := e.Group("/admin")
	g.Use(middleware.RequireFirebaseAuth(authClient))
	g.Use(middleware.RequireAdminRole(adminStore))

	g.GET("/", ctrl.AdminCheck)
	g.GET("/contests/list", ctrl.ListContests)
	g.POST("/contest", ctrl.CreateContest)
	g.PUT("/contest/:id", ctrl.UpdateContest)
	g.DELETE("/contest/:id", ctrl.DeleteContest)
	g.POST("/contest/:id/finalize", ctrl.FinalizeContest)
	g.POST("/:contestId/problem", ctrl.AddProblem)
	g.PUT("/:contestId/:problemId", ctrl.UpdateProblem)
	g.DELETE("/:contestId/:problemId", ctrl.DeleteProblem)
	g.PUT("/:contestId/leaderboard/:userId", ctrl.UpdateLeaderboardFlags)
	g.GET("/contests/:contestId/registrations", ctrl.ListRegistrations)
}
