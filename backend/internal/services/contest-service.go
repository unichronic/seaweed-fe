package services

import (
	"context"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
	"github.com/google/uuid"
)

type ContestService struct {
	store *stores.ContestStore
}

func NewContestService(store *stores.ContestStore) *ContestService {
	return &ContestService{store: store}
}

func (s *ContestService) CreateContest(ctx context.Context, c *models.Contest) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return s.store.CreateContest(ctx, c)
}

func (s *ContestService) ListContests(ctx context.Context) ([]models.Contest, error) {
	return s.store.ListContests(ctx)
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
