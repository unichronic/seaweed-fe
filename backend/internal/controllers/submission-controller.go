package controllers

import (
	"net/http"

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

type SubmissionRequest struct {
	ContestID string `json:"contest_id"`
	ProblemID string `json:"problem_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}

func (ctrl *SubmissionController) Submit(c echo.Context) error {
	userId := c.Get("userId").(string)
	req := new(SubmissionRequest)
	if err := c.Bind(req); err != nil {
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to submit")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"submission_id": submissionId})
}
