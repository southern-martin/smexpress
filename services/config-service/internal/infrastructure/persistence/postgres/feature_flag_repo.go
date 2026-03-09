package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
)

type FeatureFlagRepo struct {
	pool *pgxpool.Pool
}

func NewFeatureFlagRepo(pool *pgxpool.Pool) *FeatureFlagRepo {
	return &FeatureFlagRepo{pool: pool}
}

func (r *FeatureFlagRepo) Create(ctx context.Context, flag *entity.FeatureFlag) error {
	query := `INSERT INTO feature_flags (country_code, flag_key, enabled, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		flag.CountryCode, flag.FlagKey, flag.Enabled, flag.Description,
	).Scan(&flag.ID, &flag.CreatedAt, &flag.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: flag already exists for country", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert feature_flag: %w", err)
	}
	return nil
}

func (r *FeatureFlagRepo) GetByID(ctx context.Context, id string) (*entity.FeatureFlag, error) {
	query := `SELECT id, country_code, flag_key, enabled, description, created_at, updated_at
		FROM feature_flags WHERE id = $1`
	f := &entity.FeatureFlag{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.CountryCode, &f.FlagKey, &f.Enabled, &f.Description, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("get feature_flag: %w", err)
	}
	return f, nil
}

func (r *FeatureFlagRepo) Update(ctx context.Context, flag *entity.FeatureFlag) error {
	query := `UPDATE feature_flags SET flag_key=$1, enabled=$2, description=$3, updated_at=NOW()
		WHERE id = $4 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query, flag.FlagKey, flag.Enabled, flag.Description, flag.ID).Scan(&flag.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update feature_flag: %w", err)
	}
	return nil
}

func (r *FeatureFlagRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM feature_flags WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete feature_flag: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *FeatureFlagRepo) List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.FeatureFlag], error) {
	var where string
	var args []any

	if countryCode != "" {
		where = "WHERE country_code = $1"
		args = append(args, countryCode)
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM feature_flags %s", where)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return db.PagedResult[entity.FeatureFlag]{}, fmt.Errorf("count feature_flags: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(len(args) + 1)
	args = append(args, limitArgs...)

	dataQuery := fmt.Sprintf(`SELECT id, country_code, flag_key, enabled, description, created_at, updated_at
		FROM feature_flags %s ORDER BY flag_key %s`, where, limitClause)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return db.PagedResult[entity.FeatureFlag]{}, fmt.Errorf("list feature_flags: %w", err)
	}
	defer rows.Close()

	var items []entity.FeatureFlag
	for rows.Next() {
		var f entity.FeatureFlag
		if err := rows.Scan(&f.ID, &f.CountryCode, &f.FlagKey, &f.Enabled, &f.Description, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return db.PagedResult[entity.FeatureFlag]{}, fmt.Errorf("scan feature_flag: %w", err)
		}
		items = append(items, f)
	}

	return db.NewPagedResult(items, total, page), nil
}

func (r *FeatureFlagRepo) IsEnabled(ctx context.Context, countryCode, flagKey string) (bool, error) {
	var enabled bool
	err := r.pool.QueryRow(ctx,
		`SELECT enabled FROM feature_flags WHERE country_code = $1 AND flag_key = $2`,
		countryCode, flagKey,
	).Scan(&enabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check feature_flag: %w", err)
	}
	return enabled, nil
}
