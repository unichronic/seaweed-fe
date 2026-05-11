package stores

import (
	"context"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{pool: pool}
}

func (s *UserStore) Create(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (id, name, email, usn, mobile_number, joining_year, department)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.pool.Exec(ctx, query, u.ID, u.Name, u.Email, u.USN, u.MobileNumber, u.JoiningYear, u.Department)
	return err
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, name, email, usn, COALESCE(mobile_number, ''), joining_year, department FROM users WHERE id = $1`
	u := &models.User{}
	err := s.pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.USN, &u.MobileNumber, &u.JoiningYear, &u.Department)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserStore) Update(ctx context.Context, id string, req *models.UpdateUserProfileRequest) (*models.User, error) {
	current, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		current.Name = req.Name
	}
	if req.MobileNumber != "" {
		current.MobileNumber = req.MobileNumber
	}
	if req.Department != "" {
		current.Department = req.Department
	}
	if req.JoiningYear > 0 {
		current.JoiningYear = req.JoiningYear
	}
	query := `
		UPDATE users
		SET name = $2, mobile_number = $3, joining_year = $4, department = $5
		WHERE id = $1
	`
	_, err = s.pool.Exec(ctx, query, id, current.Name, current.MobileNumber, current.JoiningYear, current.Department)
	if err != nil {
		return nil, err
	}
	return current, nil
}
