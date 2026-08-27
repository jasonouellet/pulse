package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating pgx pool: %w", err)
	}

	if err := pool.Ping(connCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to reach PostgreSQL (ping failed): %w", err)
	}

	return pool, nil
}
