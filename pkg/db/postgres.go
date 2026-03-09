package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds PostgreSQL connection configuration.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Schema   string
	MaxConns int32
	MinConns int32
}

// DSN returns the PostgreSQL connection string.
func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Schema,
	)
}

// NewPool creates a new PostgreSQL connection pool.
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	} else {
		poolCfg.MaxConns = 10
	}
	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	} else {
		poolCfg.MinConns = 2
	}
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}
