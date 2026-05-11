package controllers

import (
	"net/http"

	"github.com/AniketSrivastava1/recruit/backend/internal/middleware"
	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/AniketSrivastava1/recruit/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type SubmissionController struct {
	service *services.SubmissionService
}

func NewSubmissionController(service *services.SubmissionService) *SubmissionController {
	return &SubmissionController{service: service}
}

func (ctrl *SubmissionController) Submit(c echo.Context) error {
	userId := c.Get("userId").(string)
	req := new(models.SubmissionRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := middleware.ValidateRequest(req); err != nil {
		return err
	}

	sub := &models.Submission{
		UserID:    userId,
		ContestID: req.ContestID,
		ProblemID: req.ProblemID,
		Language:  req.Language,
	}

	submissionId, err := ctrl.service.Submit(c.Request().Context(), sub, req.Code)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, map[string]string{"submission_id": submissionId})
}

func (ctrl *SubmissionController) Status(c echo.Context) error {
	sub, err := ctrl.service.GetStatus(c.Request().Context(), c.Param("id"), c.Get("userId").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Submission not found")
	}
	return c.JSON(http.StatusOK, sub)
}

func (ctrl *SubmissionController) Details(c echo.Context) error {
	details, err := ctrl.service.GetDetails(c.Request().Context(), c.Param("id"), c.Get("userId").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Submission not found")
	}
	return c.JSON(http.StatusOK, details)
}

func (ctrl *SubmissionController) List(c echo.Context) error {
	problemId := c.QueryParam("problem_id")
	if problemId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "problem_id is required")
	}
	submissions, err := ctrl.service.List(c.Request().Context(), c.Get("userId").(string), problemId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list submissions")
	}
	return c.JSON(http.StatusOK, submissions)
}
