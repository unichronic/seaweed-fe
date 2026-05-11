package routes

import (
	"firebase.google.com/go/auth"
	"github.com/AniketSrivastava1/recruit/backend/internal/controllers"
	"github.com/AniketSrivastava1/recruit/backend/internal/middleware"
	"github.com/labstack/echo/v4"
)

func AddUserRoutes(e *echo.Echo, authClient *auth.Client, ctrl *controllers.UserController) {
	g := e.Group("/users")
	g.Use(middleware.RequireFirebaseAuth(authClient))

	g.POST("/create", ctrl.CreateUser)
	g.GET("/profile", ctrl.GetProfile)
	g.POST("/profile", ctrl.UpdateProfile)
}
