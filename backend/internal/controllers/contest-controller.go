package controllers

import (
	"net/http"

	"github.com/AniketSrivastava1/recruit/backend/internal/middleware"
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
	req := new(models.CreateContestRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := middleware.ValidateRequest(req); err != nil {
		return err
	}
	cont := contestFromRequest(req)
	if err := ctrl.service.CreateContest(c.Request().Context(), cont); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, cont)
}

func (ctrl *ContestController) UpdateContest(c echo.Context) error {
	req := new(models.CreateContestRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := middleware.ValidateRequest(req); err != nil {
		return err
	}
	cont := contestFromRequest(req)
	cont.ID = c.Param("id")
	if err := ctrl.service.UpdateContest(c.Request().Context(), cont); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cont)
}

func (ctrl *ContestController) DeleteContest(c echo.Context) error {
	if err := ctrl.service.DeleteContest(c.Request().Context(), c.Param("id")); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete contest")
	}
	return c.NoContent(http.StatusNoContent)
}

func (ctrl *ContestController) ListContests(c echo.Context) error {
	contests, err := ctrl.service.ListContests(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list contests")
	}
	return c.JSON(http.StatusOK, contests)
}

func (ctrl *ContestController) GetContest(c echo.Context) error {
	contest, err := ctrl.service.GetContest(c.Request().Context(), c.Param("id"), optionalUserID(c))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Contest not found")
	}
	return c.JSON(http.StatusOK, contest)
}

func (ctrl *ContestController) Register(c echo.Context) error {
	req := new(models.ContestRegistrationRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := middleware.ValidateRequest(req); err != nil {
		return err
	}
	if err := ctrl.service.Register(c.Request().Context(), c.Param("id"), c.Get("userId").(string), req.Action); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"status": req.Action + "ed"})
}

func (ctrl *ContestController) AddProblem(c echo.Context) error {
	contestId := c.Param("contestId")
	req := new(models.ProblemRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := middleware.ValidateRequest(req); err != nil {
		return err
	}
	p := problemFromRequest(req, contestId, "")
	if err := ctrl.service.AddProblem(c.Request().Context(), p); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add problem")
	}
	return c.JSON(http.StatusCreated, p)
}

func (ctrl *ContestController) UpdateProblem(c echo.Context) error {
	req := new(models.ProblemRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := middleware.ValidateRequest(req); err != nil {
		return err
	}
	p := problemFromRequest(req, c.Param("contestId"), c.Param("problemId"))
	if err := ctrl.service.UpdateProblem(c.Request().Context(), p); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update problem")
	}
	return c.JSON(http.StatusOK, p)
}

func (ctrl *ContestController) DeleteProblem(c echo.Context) error {
	if err := ctrl.service.DeleteProblem(c.Request().Context(), c.Param("contestId"), c.Param("problemId")); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete problem")
	}
	return c.NoContent(http.StatusNoContent)
}

func (ctrl *ContestController) ListProblems(c echo.Context) error {
	problems, err := ctrl.service.ListProblemsForUser(c.Request().Context(), c.Param("id"), c.Get("userId").(string))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, problems)
}

func (ctrl *ContestController) GetProblem(c echo.Context) error {
	problem, err := ctrl.service.GetProblemForUser(c.Request().Context(), c.Param("id"), c.Param("problemId"), c.Get("userId").(string))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, problem)
}

func (ctrl *ContestController) GetLeaderboard(c echo.Context) error {
	rankings, err := ctrl.service.GetLeaderboard(c.Request().Context(), c.Param("id"), optionalUserID(c))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load leaderboard")
	}
	return c.JSON(http.StatusOK, rankings)
}

func (ctrl *ContestController) UpdateLeaderboardFlags(c echo.Context) error {
	req := new(models.LeaderboardUpdateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := ctrl.service.UpdateLeaderboardFlags(c.Request().Context(), c.Param("contestId"), c.Param("userId"), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update leaderboard")
	}
	return c.NoContent(http.StatusNoContent)
}

func (ctrl *ContestController) ListRegistrations(c echo.Context) error {
	users, err := ctrl.service.ListRegistrations(c.Request().Context(), c.Param("contestId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list registrations")
	}
	return c.JSON(http.StatusOK, users)
}

func (ctrl *ContestController) FinalizeContest(c echo.Context) error {
	if err := ctrl.service.FinalizeContest(c.Request().Context(), c.Param("id")); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to finalize contest")
	}
	return c.NoContent(http.StatusNoContent)
}

func (ctrl *ContestController) AdminCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{"admin": true})
}

func contestFromRequest(req *models.CreateContestRequest) *models.Contest {
	return &models.Contest{
		Name:                  req.Name,
		Description:           req.Description,
		EligibleTo:            req.EligibleTo,
		RegistrationStatus:    req.RegistrationStatus,
		RegistrationStartTime: req.RegistrationStartTime,
		RegistrationEndTime:   req.RegistrationEndTime,
		StartTime:             req.StartTime,
		EndTime:               req.EndTime,
	}
}

func problemFromRequest(req *models.ProblemRequest, contestId, problemId string) *models.Problem {
	return &models.Problem{
		ID:          problemId,
		ContestID:   contestId,
		Name:        req.Name,
		Description: req.Description,
		Score:       req.Score,
		TestCases:   req.TestCases,
	}
}

func optionalUserID(c echo.Context) string {
	value := c.Get("userId")
	userId, ok := value.(string)
	if !ok {
		return ""
	}
	return userId
}
