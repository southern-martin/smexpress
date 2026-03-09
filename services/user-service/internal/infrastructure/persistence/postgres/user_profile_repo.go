package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/user-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/user-service/internal/domain/errors"
)

type UserProfileRepo struct {
	pool *pgxpool.Pool
}

func NewUserProfileRepo(pool *pgxpool.Pool) *UserProfileRepo {
	return &UserProfileRepo{pool: pool}
}

func (r *UserProfileRepo) Create(ctx context.Context, p *entity.UserProfile) error {
	query := `INSERT INTO user_profiles (user_id, country_code, phone, mobile, job_title, department, avatar_url, timezone, locale)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		p.UserID, p.CountryCode, p.Phone, p.Mobile, p.JobTitle, p.Department, p.AvatarURL, p.Timezone, p.Locale,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("%w: profile already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert profile: %w", err)
	}
	return nil
}

func (r *UserProfileRepo) GetByUserID(ctx context.Context, userID string) (*entity.UserProfile, error) {
	p := &entity.UserProfile{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, country_code, phone, mobile, job_title, department, avatar_url, timezone, locale, created_at, updated_at
		FROM user_profiles WHERE user_id = $1`, userID,
	).Scan(&p.ID, &p.UserID, &p.CountryCode, &p.Phone, &p.Mobile, &p.JobTitle,
		&p.Department, &p.AvatarURL, &p.Timezone, &p.Locale, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return p, nil
}

func (r *UserProfileRepo) Update(ctx context.Context, p *entity.UserProfile) error {
	query := `UPDATE user_profiles SET phone=$1, mobile=$2, job_title=$3, department=$4,
		avatar_url=$5, timezone=$6, locale=$7, updated_at=NOW()
		WHERE user_id = $8 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query,
		p.Phone, p.Mobile, p.JobTitle, p.Department, p.AvatarURL, p.Timezone, p.Locale, p.UserID,
	).Scan(&p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update profile: %w", err)
	}
	return nil
}

func (r *UserProfileRepo) List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.UserProfile], error) {
	where := ""
	var args []any
	if countryCode != "" {
		where = "WHERE country_code = $1"
		args = append(args, countryCode)
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM user_profiles "+where, args...).Scan(&total); err != nil {
		return db.PagedResult[entity.UserProfile]{}, fmt.Errorf("count profiles: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(len(args) + 1)
	args = append(args, limitArgs...)

	query := fmt.Sprintf(`SELECT id, user_id, country_code, phone, mobile, job_title, department, avatar_url, timezone, locale, created_at, updated_at
		FROM user_profiles %s ORDER BY created_at DESC %s`, where, limitClause)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return db.PagedResult[entity.UserProfile]{}, fmt.Errorf("list profiles: %w", err)
	}
	defer rows.Close()

	var items []entity.UserProfile
	for rows.Next() {
		var p entity.UserProfile
		if err := rows.Scan(&p.ID, &p.UserID, &p.CountryCode, &p.Phone, &p.Mobile, &p.JobTitle,
			&p.Department, &p.AvatarURL, &p.Timezone, &p.Locale, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return db.PagedResult[entity.UserProfile]{}, fmt.Errorf("scan profile: %w", err)
		}
		items = append(items, p)
	}

	return db.NewPagedResult(items, total, page), nil
}
