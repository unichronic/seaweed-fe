package db

import (
	"context"
	"os"
	"strconv"
	"time"

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
	if maxOpen := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpen != "" {
		value, err := strconv.ParseInt(maxOpen, 10, 32)
		if err != nil {
			return nil, err
		}
		config.MaxConns = int32(value)
	}
	if maxIdle := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdle != "" {
		value, err := strconv.ParseInt(maxIdle, 10, 32)
		if err != nil {
			return nil, err
		}
		config.MinConns = int32(value)
	}
	if lifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); lifetime != "" {
		value, err := time.ParseDuration(lifetime)
		if err != nil {
			return nil, err
		}
		config.MaxConnLifetime = value
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
