package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewDBConn(lifecycle fx.Lifecycle, logger *zap.Logger) (*pgxpool.Pool, error) {
	dbAddr := os.Getenv("DB_ADDR")
	if dbAddr == "" {
		logger.Fatal("DB_ADDR environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(dbAddr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection pool")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
