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
	files := []string{
		"migrations/001_create_tables.up.sql",
		"migrations/002_seed_data.up.sql",
	}
	for _, f := range files {
		data, err := MigrationsFS.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}
		if _, err := pool.Exec(ctx, string(data)); err != nil {
			logger.Warn("migration skipped (already applied or error)", slog.String("file", f), slog.String("error", err.Error()))
		}
	}
	return nil
}
