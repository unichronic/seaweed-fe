package controllers

import (
	"net/http"

	"github.com/AniketSrivastava1/recruit/backend/internal/middleware"
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
	req := new(models.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := middleware.ValidateRequest(req); err != nil {
		return err
	}
	u := &models.User{
		ID:           userId,
		Name:         req.Name,
		Email:        req.Email,
		USN:          req.USN,
		MobileNumber: req.MobileNumber,
		JoiningYear:  req.JoiningYear,
		Department:   req.Department,
	}

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

func (ctrl *UserController) UpdateProfile(c echo.Context) error {
	userId := c.Get("userId").(string)
	req := new(models.UpdateUserProfileRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := middleware.ValidateRequest(req); err != nil {
		return err
	}
	u, err := ctrl.service.UpdateProfile(c.Request().Context(), userId, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update profile")
	}
	return c.JSON(http.StatusOK, u)
}
