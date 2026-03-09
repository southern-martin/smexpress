package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (country_code, email, password_hash, first_name, last_name, franchise_id)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''))
		RETURNING id, is_active, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		user.CountryCode, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.FranchiseID,
	).Scan(&user.ID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: email already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return r.scanUser(ctx, `SELECT id, country_code, email, password_hash, first_name, last_name,
		is_active, is_locked, failed_login_attempts, last_login_at, password_changed_at,
		COALESCE(franchise_id::text, ''), created_at, updated_at
		FROM users WHERE id = $1`, id)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.scanUser(ctx, `SELECT id, country_code, email, password_hash, first_name, last_name,
		is_active, is_locked, failed_login_attempts, last_login_at, password_changed_at,
		COALESCE(franchise_id::text, ''), created_at, updated_at
		FROM users WHERE email = $1`, email)
}

func (r *UserRepo) scanUser(ctx context.Context, query string, args ...any) (*entity.User, error) {
	u := &entity.User{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.CountryCode, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&u.IsActive, &u.IsLocked, &u.FailedLoginAttempts, &u.LastLoginAt, &u.PasswordChangedAt,
		&u.FranchiseID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query user: %w", err)
	}
	return u, nil
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET first_name=$1, last_name=$2, is_active=$3, is_locked=$4,
		failed_login_attempts=$5, last_login_at=$6, password_hash=$7, password_changed_at=$8, updated_at=NOW()
		WHERE id = $9 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query,
		user.FirstName, user.LastName, user.IsActive, user.IsLocked,
		user.FailedLoginAttempts, user.LastLoginAt, user.PasswordHash, user.PasswordChangedAt, user.ID,
	).Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *UserRepo) List(ctx context.Context, countryCode, search string, page db.Page) (db.PagedResult[entity.User], error) {
	where := "WHERE 1=1"
	args := []any{}
	argIdx := 0

	if countryCode != "" {
		argIdx++
		where += fmt.Sprintf(" AND country_code = $%d", argIdx)
		args = append(args, countryCode)
	}
	if search != "" {
		argIdx++
		where += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users "+where, args...).Scan(&total); err != nil {
		return db.PagedResult[entity.User]{}, fmt.Errorf("count users: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(argIdx + 1)
	args = append(args, limitArgs...)

	query := fmt.Sprintf(`SELECT id, country_code, email, password_hash, first_name, last_name,
		is_active, is_locked, failed_login_attempts, last_login_at, password_changed_at,
		COALESCE(franchise_id::text, ''), created_at, updated_at
		FROM users %s ORDER BY created_at DESC %s`, where, limitClause)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return db.PagedResult[entity.User]{}, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var items []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.CountryCode, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
			&u.IsActive, &u.IsLocked, &u.FailedLoginAttempts, &u.LastLoginAt, &u.PasswordChangedAt,
			&u.FranchiseID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return db.PagedResult[entity.User]{}, fmt.Errorf("scan user: %w", err)
		}
		items = append(items, u)
	}

	return db.NewPagedResult(items, total, page), nil
}
