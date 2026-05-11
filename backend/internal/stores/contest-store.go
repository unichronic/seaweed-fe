package stores

import (
	"context"
	"encoding/json"
	"time"

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
		INSERT INTO contests (id, name, description, eligible_to, registration_status, registration_start_time, registration_end_time, start_time, end_time, finalized)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := s.pool.Exec(ctx, query, c.ID, c.Name, c.Description, c.EligibleTo, c.RegistrationStatus, c.RegistrationStartTime, c.RegistrationEndTime, c.StartTime, c.EndTime, c.Finalized)
	return err
}

func (s *ContestStore) UpdateContest(ctx context.Context, c *models.Contest) error {
	query := `
		UPDATE contests
		SET name = $2, description = $3, eligible_to = $4, registration_status = $5,
			registration_start_time = $6, registration_end_time = $7, start_time = $8, end_time = $9
		WHERE id = $1
	`
	_, err := s.pool.Exec(ctx, query, c.ID, c.Name, c.Description, c.EligibleTo, c.RegistrationStatus, c.RegistrationStartTime, c.RegistrationEndTime, c.StartTime, c.EndTime)
	return err
}

func (s *ContestStore) DeleteContest(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM contests WHERE id = $1`, id)
	return err
}

func (s *ContestStore) ListContests(ctx context.Context) ([]models.Contest, error) {
	query := `SELECT id, name, description, eligible_to, registration_status, registration_start_time, registration_end_time, start_time, end_time, finalized FROM contests ORDER BY start_time DESC`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contests []models.Contest
	for rows.Next() {
		var c models.Contest
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.EligibleTo, &c.RegistrationStatus, &c.RegistrationStartTime, &c.RegistrationEndTime, &c.StartTime, &c.EndTime, &c.Finalized); err != nil {
			return nil, err
		}
		contests = append(contests, c)
	}
	return contests, nil
}

func (s *ContestStore) GetContest(ctx context.Context, id string) (*models.Contest, error) {
	query := `SELECT id, name, description, eligible_to, registration_status, registration_start_time, registration_end_time, start_time, end_time, finalized FROM contests WHERE id = $1`
	c := &models.Contest{}
	err := s.pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Description, &c.EligibleTo, &c.RegistrationStatus, &c.RegistrationStartTime, &c.RegistrationEndTime, &c.StartTime, &c.EndTime, &c.Finalized)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ContestStore) IsRegistered(ctx context.Context, contestId, userId string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM contest_registrations WHERE contest_id = $1 AND user_id = $2)`
	var exists bool
	err := s.pool.QueryRow(ctx, query, contestId, userId).Scan(&exists)
	return exists, err
}

func (s *ContestStore) Register(ctx context.Context, contestId, userId string) error {
	query := `
		INSERT INTO contest_registrations (contest_id, user_id, registered_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (contest_id, user_id) DO NOTHING
	`
	_, err := s.pool.Exec(ctx, query, contestId, userId, time.Now().UnixMilli())
	return err
}

func (s *ContestStore) Unregister(ctx context.Context, contestId, userId string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM contest_registrations WHERE contest_id = $1 AND user_id = $2`, contestId, userId)
	return err
}

func (s *ContestStore) ListRegistrations(ctx context.Context, contestId string) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.usn, COALESCE(u.mobile_number, ''), u.joining_year, u.department
		FROM contest_registrations cr
		JOIN users u ON u.id = cr.user_id
		WHERE cr.contest_id = $1
		ORDER BY cr.registered_at ASC
	`
	rows, err := s.pool.Query(ctx, query, contestId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.USN, &u.MobileNumber, &u.JoiningYear, &u.Department); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *ContestStore) AddProblem(ctx context.Context, p *models.Problem) error {
	testCases, err := json.Marshal(p.TestCases)
	if err != nil {
		return err
	}
	query := `INSERT INTO problems (id, contest_id, name, description, score, test_cases) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = s.pool.Exec(ctx, query, p.ID, p.ContestID, p.Name, p.Description, p.Score, testCases)
	return err
}

func (s *ContestStore) ListProblems(ctx context.Context, contestId string) ([]models.Problem, error) {
	query := `SELECT id, contest_id, name, description, score FROM problems WHERE contest_id = $1 ORDER BY name`
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
	return problems, rows.Err()
}

func (s *ContestStore) GetProblem(ctx context.Context, contestId, problemId string, includeTestCases bool) (*models.Problem, error) {
	p := &models.Problem{}
	if includeTestCases {
		var raw []byte
		query := `SELECT id, contest_id, name, description, score, test_cases FROM problems WHERE contest_id = $1 AND id = $2`
		err := s.pool.QueryRow(ctx, query, contestId, problemId).Scan(&p.ID, &p.ContestID, &p.Name, &p.Description, &p.Score, &raw)
		if err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &p.TestCases); err != nil {
				return nil, err
			}
		}
		return p, nil
	}
	query := `SELECT id, contest_id, name, description, score FROM problems WHERE contest_id = $1 AND id = $2`
	err := s.pool.QueryRow(ctx, query, contestId, problemId).Scan(&p.ID, &p.ContestID, &p.Name, &p.Description, &p.Score)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ContestStore) UpdateProblem(ctx context.Context, p *models.Problem) error {
	testCases, err := json.Marshal(p.TestCases)
	if err != nil {
		return err
	}
	query := `
		UPDATE problems
		SET name = $3, description = $4, score = $5, test_cases = $6
		WHERE contest_id = $1 AND id = $2
	`
	_, err = s.pool.Exec(ctx, query, p.ContestID, p.ID, p.Name, p.Description, p.Score, testCases)
	return err
}

func (s *ContestStore) DeleteProblem(ctx context.Context, contestId, problemId string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM problems WHERE contest_id = $1 AND id = $2`, contestId, problemId)
	return err
}

func (s *ContestStore) RefreshRankingMV(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY ranking_mv`)
	return err
}

func (s *ContestStore) FinalizeContest(ctx context.Context, contestId string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE contests SET finalized = TRUE WHERE id = $1`, contestId); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *ContestStore) GetLeaderboard(ctx context.Context, contestId string, includeHidden bool) ([]models.Ranking, error) {
	query := `
		SELECT contest_id, user_id, name, usn, department, score, hidden, disqualified, shortlisted,
			correct_attempts, incorrect_attempts, rank
		FROM ranking_mv
		WHERE contest_id = $1
	`
	if !includeHidden {
		query += ` AND hidden = FALSE`
	}
	query += ` ORDER BY rank, score DESC, name`

	rows, err := s.pool.Query(ctx, query, contestId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rankings := []models.Ranking{}
	for rows.Next() {
		var r models.Ranking
		if err := rows.Scan(&r.ContestID, &r.UserID, &r.Name, &r.USN, &r.Department, &r.Score, &r.Hidden, &r.Disqualified, &r.Shortlisted, &r.CorrectAttempts, &r.IncorrectAttempts, &r.Rank); err != nil {
			return nil, err
		}
		rankings = append(rankings, r)
	}
	return rankings, rows.Err()
}

func (s *ContestStore) UpdateLeaderboardFlags(ctx context.Context, contestId, userId string, hidden, disqualified, shortlisted *bool) error {
	query := `
		INSERT INTO rankings (contest_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (contest_id, user_id) DO NOTHING
	`
	if _, err := s.pool.Exec(ctx, query, contestId, userId); err != nil {
		return err
	}
	if hidden != nil {
		if _, err := s.pool.Exec(ctx, `UPDATE rankings SET hidden = $3 WHERE contest_id = $1 AND user_id = $2`, contestId, userId, *hidden); err != nil {
			return err
		}
	}
	if disqualified != nil {
		if _, err := s.pool.Exec(ctx, `UPDATE rankings SET disqualified = $3 WHERE contest_id = $1 AND user_id = $2`, contestId, userId, *disqualified); err != nil {
			return err
		}
	}
	if shortlisted != nil {
		if _, err := s.pool.Exec(ctx, `UPDATE rankings SET shortlisted = $3 WHERE contest_id = $1 AND user_id = $2`, contestId, userId, *shortlisted); err != nil {
			return err
		}
	}
	return s.RefreshRankingMV(ctx)
}
