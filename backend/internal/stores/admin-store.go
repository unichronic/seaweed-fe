package stores

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminStore struct {
	pool *pgxpool.Pool
}

func NewAdminStore(pool *pgxpool.Pool) *AdminStore {
	return &AdminStore{pool: pool}
}

func (s *AdminStore) IsAdmin(ctx context.Context, userId string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM admin WHERE user_id = $1)`
	var exists bool
	err := s.pool.QueryRow(ctx, query, userId).Scan(&exists)
	return exists, err
}
