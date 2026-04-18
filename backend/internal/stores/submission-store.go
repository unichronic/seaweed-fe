package stores

import (
	"context"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubmissionStore struct {
	pool *pgxpool.Pool
}

func NewSubmissionStore(pool *pgxpool.Pool) *SubmissionStore {
	return &SubmissionStore{pool: pool}
}

func (s *SubmissionStore) CreateSubmission(ctx context.Context, sub *models.Submission) error {
	query := `
		INSERT INTO submissions (id, user_id, contest_id, problem_id, language, s3_key, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.pool.Exec(ctx, query, sub.ID, sub.UserID, sub.ContestID, sub.ProblemID, sub.Language, sub.S3Key, sub.Status, sub.CreatedAt)
	return err
}

func (s *SubmissionStore) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE submissions SET status = $1 WHERE id = $2`
	_, err := s.pool.Exec(ctx, query, status, id)
	return err
}

func (s *SubmissionStore) SaveTestCaseResult(ctx context.Context, tc *models.TestCaseResult) error {
	query := `
		INSERT INTO test_case_results (submission_id, test_case_id, status, runtime_ms, memory_kb, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.pool.Exec(ctx, query, tc.SubmissionID, tc.TestCaseID, tc.Status, tc.RuntimeMS, tc.MemoryKB, tc.CreatedAt)
	return err
}

func (s *SubmissionStore) RefreshRankingMV(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY ranking_mv")
	return err
}

func (s *SubmissionStore) UpdateRanking(ctx context.Context, contestId, userId string, points int) error {
	query := `
		INSERT INTO rankings (contest_id, user_id, score)
		VALUES ($1, $2, $3)
		ON CONFLICT (contest_id, user_id)
		DO UPDATE SET score = rankings.score + EXCLUDED.score
	`
	_, err := s.pool.Exec(ctx, query, contestId, userId, points)
	return err
}
