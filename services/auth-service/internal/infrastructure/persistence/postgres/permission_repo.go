package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
)

type PermissionRepo struct {
	pool *pgxpool.Pool
}

func NewPermissionRepo(pool *pgxpool.Pool) *PermissionRepo {
	return &PermissionRepo{pool: pool}
}

func (r *PermissionRepo) List(ctx context.Context) ([]entity.Permission, error) {
	return r.queryPermissions(ctx, `SELECT id, code, name, module, COALESCE(description, ''), created_at FROM permissions ORDER BY module, code`)
}

func (r *PermissionRepo) ListByModule(ctx context.Context, module string) ([]entity.Permission, error) {
	return r.queryPermissions(ctx, `SELECT id, code, name, module, COALESCE(description, ''), created_at FROM permissions WHERE module = $1 ORDER BY code`, module)
}

func (r *PermissionRepo) queryPermissions(ctx context.Context, query string, args ...any) ([]entity.Permission, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query permissions: %w", err)
	}
	defer rows.Close()

	var items []entity.Permission
	for rows.Next() {
		var p entity.Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.Module, &p.Description, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		items = append(items, p)
	}
	return items, nil
}
