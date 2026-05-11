package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ContestService struct {
	store      *stores.ContestStore
	adminStore *stores.AdminStore
}

func NewContestService(store *stores.ContestStore, adminStore *stores.AdminStore) *ContestService {
	return &ContestService{store: store, adminStore: adminStore}
}

func (s *ContestService) CreateContest(ctx context.Context, c *models.Contest) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.RegistrationStatus == "" {
		c.RegistrationStatus = "open"
	}
	if err := validateContestTimes(c); err != nil {
		return err
	}
	return s.store.CreateContest(ctx, c)
}

func (s *ContestService) UpdateContest(ctx context.Context, c *models.Contest) error {
	if c.RegistrationStatus == "" {
		c.RegistrationStatus = "open"
	}
	if err := validateContestTimes(c); err != nil {
		return err
	}
	return s.store.UpdateContest(ctx, c)
}

func (s *ContestService) DeleteContest(ctx context.Context, id string) error {
	return s.store.DeleteContest(ctx, id)
}

func (s *ContestService) ListContests(ctx context.Context) ([]models.Contest, error) {
	return s.store.ListContests(ctx)
}

func (s *ContestService) GetContest(ctx context.Context, id string, userId string) (*models.ContestResponse, error) {
	contest, err := s.store.GetContest(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := &models.ContestResponse{Contest: *contest}
	if userId != "" {
		registered, err := s.store.IsRegistered(ctx, id, userId)
		if err != nil {
			return nil, err
		}
		resp.Registered = registered
	}
	return resp, nil
}

func (s *ContestService) Register(ctx context.Context, contestId, userId, action string) error {
	contest, err := s.store.GetContest(ctx, contestId)
	if err != nil {
		return err
	}
	if action == "unregister" {
		if time.Now().UnixMilli() >= contest.StartTime {
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot unregister after contest start")
		}
		return s.store.Unregister(ctx, contestId, userId)
	}
	now := time.Now().UnixMilli()
	if contest.RegistrationStatus != "open" {
		return echo.NewHTTPError(http.StatusForbidden, "Contest registration is not open")
	}
	if now < contest.RegistrationStartTime || now > contest.RegistrationEndTime {
		return echo.NewHTTPError(http.StatusForbidden, "Contest registration window is closed")
	}
	return s.store.Register(ctx, contestId, userId)
}

func (s *ContestService) AddProblem(ctx context.Context, p *models.Problem) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return s.store.AddProblem(ctx, p)
}

func (s *ContestService) ListProblems(ctx context.Context, contestId string) ([]models.Problem, error) {
	return s.store.ListProblems(ctx, contestId)
}

func (s *ContestService) ListProblemsForUser(ctx context.Context, contestId, userId string) ([]models.ProblemSummaryResponse, error) {
	if err := s.ensureContestAccess(ctx, contestId, userId); err != nil {
		return nil, err
	}
	problems, err := s.store.ListProblems(ctx, contestId)
	if err != nil {
		return nil, err
	}
	resp := make([]models.ProblemSummaryResponse, 0, len(problems))
	for _, p := range problems {
		resp = append(resp, models.ProblemSummaryResponse{ID: p.ID, ContestID: p.ContestID, Name: p.Name, Score: p.Score})
	}
	return resp, nil
}

func (s *ContestService) GetProblemForUser(ctx context.Context, contestId, problemId, userId string) (*models.ProblemDetailResponse, error) {
	if err := s.ensureContestAccess(ctx, contestId, userId); err != nil {
		return nil, err
	}
	p, err := s.store.GetProblem(ctx, contestId, problemId, false)
	if err != nil {
		return nil, err
	}
	return &models.ProblemDetailResponse{ID: p.ID, ContestID: p.ContestID, Name: p.Name, Description: p.Description, Score: p.Score}, nil
}

func (s *ContestService) GetProblemForJudge(ctx context.Context, contestId, problemId string) (*models.Problem, error) {
	return s.store.GetProblem(ctx, contestId, problemId, true)
}

func (s *ContestService) UpdateProblem(ctx context.Context, p *models.Problem) error {
	return s.store.UpdateProblem(ctx, p)
}

func (s *ContestService) DeleteProblem(ctx context.Context, contestId, problemId string) error {
	return s.store.DeleteProblem(ctx, contestId, problemId)
}

func (s *ContestService) GetLeaderboard(ctx context.Context, contestId, userId string) ([]models.Ranking, error) {
	includeHidden := false
	if userId != "" && s.adminStore != nil {
		isAdmin, err := s.adminStore.IsAdmin(ctx, userId)
		if err != nil {
			return nil, err
		}
		includeHidden = isAdmin
	}
	return s.store.GetLeaderboard(ctx, contestId, includeHidden)
}

func (s *ContestService) UpdateLeaderboardFlags(ctx context.Context, contestId, userId string, req *models.LeaderboardUpdateRequest) error {
	return s.store.UpdateLeaderboardFlags(ctx, contestId, userId, req.Hidden, req.Disqualified, req.Shortlisted)
}

func (s *ContestService) ListRegistrations(ctx context.Context, contestId string) ([]models.User, error) {
	return s.store.ListRegistrations(ctx, contestId)
}

func (s *ContestService) FinalizeContest(ctx context.Context, contestId string) error {
	if err := s.store.RefreshRankingMV(ctx); err != nil {
		return err
	}
	return s.store.FinalizeContest(ctx, contestId)
}

func (s *ContestService) EnsureContestAccess(ctx context.Context, contestId, userId string) error {
	return s.ensureContestAccess(ctx, contestId, userId)
}

func (s *ContestService) ensureContestAccess(ctx context.Context, contestId, userId string) error {
	if userId == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	contest, err := s.store.GetContest(ctx, contestId)
	if err != nil {
		return err
	}
	registered, err := s.store.IsRegistered(ctx, contestId, userId)
	if err != nil {
		return err
	}
	if !registered {
		return echo.NewHTTPError(http.StatusForbidden, "Register for the contest before accessing problems")
	}
	now := time.Now().UnixMilli()
	if now < contest.StartTime || now > contest.EndTime {
		return echo.NewHTTPError(http.StatusForbidden, "Problems are available only during the contest window")
	}
	return nil
}

func validateContestTimes(c *models.Contest) error {
	if c.RegistrationStartTime >= c.RegistrationEndTime {
		return echo.NewHTTPError(http.StatusBadRequest, "Registration start must be before registration end")
	}
	if c.StartTime >= c.EndTime {
		return echo.NewHTTPError(http.StatusBadRequest, "Contest start must be before contest end")
	}
	if c.RegistrationEndTime > c.EndTime {
		return echo.NewHTTPError(http.StatusBadRequest, "Registration cannot end after contest end")
	}
	if c.RegistrationStatus != "open" && c.RegistrationStatus != "closed" && c.RegistrationStatus != "invite-only" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid registration status: %s", c.RegistrationStatus))
	}
	return nil
}
