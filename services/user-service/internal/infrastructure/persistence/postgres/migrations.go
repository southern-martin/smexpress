package postgres

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) error {
	data, err := MigrationsFS.ReadFile("migrations/001_create_tables.up.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	if _, err := pool.Exec(ctx, string(data)); err != nil {
		logger.Warn("migration skipped", slog.String("error", err.Error()))
	}
	return nil
}
