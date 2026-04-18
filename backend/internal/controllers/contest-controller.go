package controllers

import (
	"net/http"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/AniketSrivastava1/recruit/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type ContestController struct {
	service *services.ContestService
}

func NewContestController(service *services.ContestService) *ContestController {
	return &ContestController{service: service}
}

func (ctrl *ContestController) CreateContest(c echo.Context) error {
	cont := new(models.Contest)
	if err := c.Bind(cont); err != nil {
		return err
	}
	if err := ctrl.service.CreateContest(c.Request().Context(), cont); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create contest")
	}
	return c.JSON(http.StatusCreated, cont)
}

func (ctrl *ContestController) ListContests(c echo.Context) error {
	contests, err := ctrl.service.ListContests(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list contests")
	}
	return c.JSON(http.StatusOK, contests)
}

func (ctrl *ContestController) AddProblem(c echo.Context) error {
	contestId := c.Param("contestId")
	p := new(models.Problem)
	if err := c.Bind(p); err != nil {
		return err
	}
	p.ContestID = contestId
	if err := ctrl.service.AddProblem(c.Request().Context(), p); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add problem")
	}
	return c.JSON(http.StatusCreated, p)
}

func (ctrl *ContestController) ListProblems(c echo.Context) error {
	contestId := c.Param("contestId")
	problems, err := ctrl.service.ListProblems(c.Request().Context(), contestId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list problems")
	}
	return c.JSON(http.StatusOK, problems)
}
