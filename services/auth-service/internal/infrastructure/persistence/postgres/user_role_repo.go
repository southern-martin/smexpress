package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
)

type UserRoleRepo struct {
	pool *pgxpool.Pool
}

func NewUserRoleRepo(pool *pgxpool.Pool) *UserRoleRepo {
	return &UserRoleRepo{pool: pool}
}

func (r *UserRoleRepo) AssignRoles(ctx context.Context, userID string, roleIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("clear roles: %w", err)
	}

	for _, rid := range roleIDs {
		_, err = tx.Exec(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, userID, rid)
		if err != nil {
			return fmt.Errorf("assign role: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *UserRoleRepo) GetUserRoles(ctx context.Context, userID string) ([]entity.Role, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.country_code, r.name, r.display_name, COALESCE(r.description, ''), r.is_system, r.created_at, r.updated_at
		FROM roles r JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 ORDER BY r.name`, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	defer rows.Close()

	var items []entity.Role
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.CountryCode, &role.Name, &role.DisplayName,
			&role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		items = append(items, role)
	}
	return items, nil
}

func (r *UserRoleRepo) GetUserPermissions(ctx context.Context, userID string) ([]entity.Permission, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT p.id, p.code, p.name, p.module, COALESCE(p.description, ''), p.created_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1 ORDER BY p.module, p.code`, userID)
	if err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
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
