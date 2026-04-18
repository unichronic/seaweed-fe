package stores

import (
	"context"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContestStore struct {
	pool *pgxpool.Pool
}

func NewContestStore(pool *pgxpool.Pool) *ContestStore {
	return &ContestStore{pool: pool}
}

func (s *ContestStore) CreateContest(ctx context.Context, c *models.Contest) error {
	query := `
		INSERT INTO contests (id, name, description, eligible_to, registration_start_time, registration_end_time, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.pool.Exec(ctx, query, c.ID, c.Name, c.Description, c.EligibleTo, c.RegistrationStartTime, c.RegistrationEndTime, c.StartTime, c.EndTime)
	return err
}

func (s *ContestStore) ListContests(ctx context.Context) ([]models.Contest, error) {
	query := `SELECT id, name, description, eligible_to, registration_start_time, registration_end_time, start_time, end_time FROM contests`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contests []models.Contest
	for rows.Next() {
		var c models.Contest
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.EligibleTo, &c.RegistrationStartTime, &c.RegistrationEndTime, &c.StartTime, &c.EndTime); err != nil {
			return nil, err
		}
		contests = append(contests, c)
	}
	return contests, nil
}

func (s *ContestStore) AddProblem(ctx context.Context, p *models.Problem) error {
	query := `INSERT INTO problems (id, contest_id, name, description, score) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.pool.Exec(ctx, query, p.ID, p.ContestID, p.Name, p.Description, p.Score)
	return err
}

func (s *ContestStore) ListProblems(ctx context.Context, contestId string) ([]models.Problem, error) {
	query := `SELECT id, contest_id, name, description, score FROM problems WHERE contest_id = $1`
	rows, err := s.pool.Query(ctx, query, contestId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []models.Problem
	for rows.Next() {
		var p models.Problem
		if err := rows.Scan(&p.ID, &p.ContestID, &p.Name, &p.Description, &p.Score); err != nil {
			return nil, err
		}
		problems = append(problems, p)
	}
	return problems, nil
}
