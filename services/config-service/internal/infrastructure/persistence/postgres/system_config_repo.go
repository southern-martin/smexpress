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

type SystemConfigRepo struct {
	pool *pgxpool.Pool
}

func NewSystemConfigRepo(pool *pgxpool.Pool) *SystemConfigRepo {
	return &SystemConfigRepo{pool: pool}
}

func (r *SystemConfigRepo) Create(ctx context.Context, cfg *entity.SystemConfig) error {
	query := `INSERT INTO system_configs (country_code, config_key, config_value, description, data_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		cfg.CountryCode, cfg.ConfigKey, cfg.ConfigValue, cfg.Description, cfg.DataType,
	).Scan(&cfg.ID, &cfg.CreatedAt, &cfg.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: config key already exists for country", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert system_config: %w", err)
	}
	return nil
}

func (r *SystemConfigRepo) GetByID(ctx context.Context, id string) (*entity.SystemConfig, error) {
	query := `SELECT id, country_code, config_key, config_value, description, data_type, created_at, updated_at
		FROM system_configs WHERE id = $1`
	cfg := &entity.SystemConfig{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&cfg.ID, &cfg.CountryCode, &cfg.ConfigKey, &cfg.ConfigValue,
		&cfg.Description, &cfg.DataType, &cfg.CreatedAt, &cfg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("get system_config: %w", err)
	}
	return cfg, nil
}

func (r *SystemConfigRepo) GetByKey(ctx context.Context, countryCode, key string) (*entity.SystemConfig, error) {
	query := `SELECT id, country_code, config_key, config_value, description, data_type, created_at, updated_at
		FROM system_configs WHERE country_code = $1 AND config_key = $2`
	cfg := &entity.SystemConfig{}
	err := r.pool.QueryRow(ctx, query, countryCode, key).Scan(
		&cfg.ID, &cfg.CountryCode, &cfg.ConfigKey, &cfg.ConfigValue,
		&cfg.Description, &cfg.DataType, &cfg.CreatedAt, &cfg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("get system_config by key: %w", err)
	}
	return cfg, nil
}

func (r *SystemConfigRepo) Update(ctx context.Context, cfg *entity.SystemConfig) error {
	query := `UPDATE system_configs SET config_value = $1, description = $2, data_type = $3, updated_at = NOW()
		WHERE id = $4 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query, cfg.ConfigValue, cfg.Description, cfg.DataType, cfg.ID).Scan(&cfg.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update system_config: %w", err)
	}
	return nil
}

func (r *SystemConfigRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM system_configs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete system_config: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *SystemConfigRepo) List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.SystemConfig], error) {
	var (
		where string
		args  []any
	)

	if countryCode != "" {
		where = "WHERE country_code = $1"
		args = append(args, countryCode)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM system_configs %s", where)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return db.PagedResult[entity.SystemConfig]{}, fmt.Errorf("count system_configs: %w", err)
	}

	argOffset := len(args)
	limitClause, limitArgs := page.LimitOffsetClause(argOffset + 1)
	args = append(args, limitArgs...)

	dataQuery := fmt.Sprintf(`SELECT id, country_code, config_key, config_value, description, data_type, created_at, updated_at
		FROM system_configs %s ORDER BY config_key %s`, where, limitClause)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return db.PagedResult[entity.SystemConfig]{}, fmt.Errorf("list system_configs: %w", err)
	}
	defer rows.Close()

	var items []entity.SystemConfig
	for rows.Next() {
		var c entity.SystemConfig
		if err := rows.Scan(&c.ID, &c.CountryCode, &c.ConfigKey, &c.ConfigValue,
			&c.Description, &c.DataType, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return db.PagedResult[entity.SystemConfig]{}, fmt.Errorf("scan system_config: %w", err)
		}
		items = append(items, c)
	}

	return db.NewPagedResult(items, total, page), nil
}
