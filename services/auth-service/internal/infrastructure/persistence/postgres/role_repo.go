package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
)

type RoleRepo struct {
	pool *pgxpool.Pool
}

func NewRoleRepo(pool *pgxpool.Pool) *RoleRepo {
	return &RoleRepo{pool: pool}
}

func (r *RoleRepo) Create(ctx context.Context, role *entity.Role) error {
	query := `INSERT INTO roles (country_code, name, display_name, description)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		role.CountryCode, role.Name, role.DisplayName, role.Description,
	).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: role already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert role: %w", err)
	}
	return nil
}

func (r *RoleRepo) GetByID(ctx context.Context, id string) (*entity.Role, error) {
	role := &entity.Role{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, country_code, name, display_name, description, is_system, created_at, updated_at
		FROM roles WHERE id = $1`, id,
	).Scan(&role.ID, &role.CountryCode, &role.Name, &role.DisplayName, &role.Description,
		&role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("get role: %w", err)
	}
	return role, nil
}

func (r *RoleRepo) Update(ctx context.Context, role *entity.Role) error {
	query := `UPDATE roles SET display_name=$1, description=$2, updated_at=NOW()
		WHERE id = $3 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query, role.DisplayName, role.Description, role.ID).Scan(&role.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update role: %w", err)
	}
	return nil
}

func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM roles WHERE id = $1 AND is_system = false`, id)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *RoleRepo) List(ctx context.Context, countryCode string) ([]entity.Role, error) {
	query := `SELECT id, country_code, name, display_name, description, is_system, created_at, updated_at FROM roles`
	var args []any
	if countryCode != "" {
		query += " WHERE country_code = $1"
		args = append(args, countryCode)
	}
	query += " ORDER BY name"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
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

func (r *RoleRepo) SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID)
	if err != nil {
		return fmt.Errorf("clear permissions: %w", err)
	}

	for _, pid := range permissionIDs {
		_, err = tx.Exec(ctx, `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)`, roleID, pid)
		if err != nil {
			return fmt.Errorf("assign permission: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *RoleRepo) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]entity.Permission, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.id, p.code, p.name, p.module, COALESCE(p.description, ''), p.created_at
		FROM permissions p JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1 ORDER BY p.module, p.code`, roleID)
	if err != nil {
		return nil, fmt.Errorf("get role permissions: %w", err)
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
