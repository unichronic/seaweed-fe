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
	query := `SELECT id, name, email, usn, mobile_number, joining_year, department FROM users WHERE id = $1`
	u := &models.User{}
	err := s.pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.USN, &u.MobileNumber, &u.JoiningYear, &u.Department)
	if err != nil {
		return nil, err
	}
	return u, nil
}
