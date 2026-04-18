package controllers

import (
	"net/http"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/AniketSrivastava1/recruit/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (ctrl *UserController) CreateUser(c echo.Context) error {
	userId := c.Get("userId").(string)
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	u.ID = userId

	if err := ctrl.service.CreateUser(c.Request().Context(), u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, u)
}

func (ctrl *UserController) GetProfile(c echo.Context) error {
	userId := c.Get("userId").(string)
	u, err := ctrl.service.GetUserByID(c.Request().Context(), userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, u)
}
