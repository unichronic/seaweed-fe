package stores

import (
	"context"
	"time"

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
	query := `UPDATE submissions SET status = $1, updated_at = $3 WHERE id = $2`
	_, err := s.pool.Exec(ctx, query, status, id, time.Now().UnixMilli())
	return err
}

func (s *SubmissionStore) GetSubmission(ctx context.Context, id, userId string) (*models.Submission, error) {
	query := `SELECT id, user_id, contest_id, problem_id, language, s3_key, status, created_at FROM submissions WHERE id = $1 AND user_id = $2`
	sub := &models.Submission{}
	err := s.pool.QueryRow(ctx, query, id, userId).Scan(&sub.ID, &sub.UserID, &sub.ContestID, &sub.ProblemID, &sub.Language, &sub.S3Key, &sub.Status, &sub.CreatedAt)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubmissionStore) ListSubmissions(ctx context.Context, userId, problemId string) ([]models.Submission, error) {
	query := `
		SELECT id, user_id, contest_id, problem_id, language, s3_key, status, created_at
		FROM submissions
		WHERE user_id = $1 AND problem_id = $2
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query, userId, problemId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := []models.Submission{}
	for rows.Next() {
		var sub models.Submission
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.ContestID, &sub.ProblemID, &sub.Language, &sub.S3Key, &sub.Status, &sub.CreatedAt); err != nil {
			return nil, err
		}
		submissions = append(submissions, sub)
	}
	return submissions, rows.Err()
}

func (s *SubmissionStore) GetTestCaseResults(ctx context.Context, submissionId string) ([]models.TestCaseResult, error) {
	query := `
		SELECT id, submission_id, test_case_id, status, runtime_ms, memory_kb, created_at
		FROM test_case_results
		WHERE submission_id = $1
		ORDER BY created_at ASC
	`
	rows, err := s.pool.Query(ctx, query, submissionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []models.TestCaseResult{}
	for rows.Next() {
		var result models.TestCaseResult
		if err := rows.Scan(&result.ID, &result.SubmissionID, &result.TestCaseID, &result.Status, &result.RuntimeMS, &result.MemoryKB, &result.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, rows.Err()
}

func (s *SubmissionStore) ApplyVerdict(ctx context.Context, sub *models.Submission, status string, problemScore int, results []models.TestCaseResult) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var alreadyAccepted bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM submissions
			WHERE contest_id = $1 AND problem_id = $2 AND user_id = $3 AND id <> $4 AND status = 'accepted'
		)
	`, sub.ContestID, sub.ProblemID, sub.UserID, sub.ID).Scan(&alreadyAccepted)
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	if _, err := tx.Exec(ctx, `UPDATE submissions SET status = $1, updated_at = $2 WHERE id = $3`, status, now, sub.ID); err != nil {
		return err
	}

	for _, result := range results {
		if result.CreatedAt == 0 {
			result.CreatedAt = now
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO test_case_results (submission_id, test_case_id, status, runtime_ms, memory_kb, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, sub.ID, result.TestCaseID, result.Status, result.RuntimeMS, result.MemoryKB, result.CreatedAt)
		if err != nil {
			return err
		}
	}

	points := 0
	correct := 0
	incorrect := 0
	if status == "accepted" {
		correct = 1
		if !alreadyAccepted {
			points = problemScore
		}
	} else if status != "failed_to_process" {
		incorrect = 1
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO rankings (contest_id, user_id, score, correct_attempts, incorrect_attempts)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (contest_id, user_id)
		DO UPDATE SET
			score = rankings.score + EXCLUDED.score,
			correct_attempts = rankings.correct_attempts + EXCLUDED.correct_attempts,
			incorrect_attempts = rankings.incorrect_attempts + EXCLUDED.incorrect_attempts
	`, sub.ContestID, sub.UserID, points, correct, incorrect)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return s.RefreshRankingMV(ctx)
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
