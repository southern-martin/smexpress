package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
)

type CountryConfigRepo struct {
	pool *pgxpool.Pool
}

func NewCountryConfigRepo(pool *pgxpool.Pool) *CountryConfigRepo {
	return &CountryConfigRepo{pool: pool}
}

func (r *CountryConfigRepo) Create(ctx context.Context, cfg *entity.CountryConfig) error {
	query := `INSERT INTO country_configs (country_code, country_name, currency_code, currency_symbol, timezone, date_format, weight_unit, dimension_unit, locale, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		cfg.CountryCode, cfg.CountryName, cfg.CurrencyCode, cfg.CurrencySymbol,
		cfg.Timezone, cfg.DateFormat, cfg.WeightUnit, cfg.DimensionUnit, cfg.Locale, cfg.IsActive,
	).Scan(&cfg.ID, &cfg.CreatedAt, &cfg.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: country already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert country_config: %w", err)
	}
	return nil
}

func (r *CountryConfigRepo) GetByCode(ctx context.Context, code string) (*entity.CountryConfig, error) {
	query := `SELECT id, country_code, country_name, currency_code, currency_symbol, timezone, date_format, weight_unit, dimension_unit, locale, is_active, created_at, updated_at
		FROM country_configs WHERE country_code = $1`
	cfg := &entity.CountryConfig{}
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&cfg.ID, &cfg.CountryCode, &cfg.CountryName, &cfg.CurrencyCode, &cfg.CurrencySymbol,
		&cfg.Timezone, &cfg.DateFormat, &cfg.WeightUnit, &cfg.DimensionUnit, &cfg.Locale,
		&cfg.IsActive, &cfg.CreatedAt, &cfg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("get country_config: %w", err)
	}
	return cfg, nil
}

func (r *CountryConfigRepo) Update(ctx context.Context, cfg *entity.CountryConfig) error {
	query := `UPDATE country_configs SET country_name=$1, currency_code=$2, currency_symbol=$3, timezone=$4,
		date_format=$5, weight_unit=$6, dimension_unit=$7, locale=$8, is_active=$9, updated_at=NOW()
		WHERE country_code = $10 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query,
		cfg.CountryName, cfg.CurrencyCode, cfg.CurrencySymbol, cfg.Timezone,
		cfg.DateFormat, cfg.WeightUnit, cfg.DimensionUnit, cfg.Locale, cfg.IsActive, cfg.CountryCode,
	).Scan(&cfg.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update country_config: %w", err)
	}
	return nil
}

func (r *CountryConfigRepo) List(ctx context.Context) ([]entity.CountryConfig, error) {
	return r.queryCountries(ctx, "SELECT id, country_code, country_name, currency_code, currency_symbol, timezone, date_format, weight_unit, dimension_unit, locale, is_active, created_at, updated_at FROM country_configs ORDER BY country_name")
}

func (r *CountryConfigRepo) ListActive(ctx context.Context) ([]entity.CountryConfig, error) {
	return r.queryCountries(ctx, "SELECT id, country_code, country_name, currency_code, currency_symbol, timezone, date_format, weight_unit, dimension_unit, locale, is_active, created_at, updated_at FROM country_configs WHERE is_active = true ORDER BY country_name")
}

func (r *CountryConfigRepo) queryCountries(ctx context.Context, query string) ([]entity.CountryConfig, error) {
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query country_configs: %w", err)
	}
	defer rows.Close()

	var items []entity.CountryConfig
	for rows.Next() {
		var c entity.CountryConfig
		if err := rows.Scan(&c.ID, &c.CountryCode, &c.CountryName, &c.CurrencyCode, &c.CurrencySymbol,
			&c.Timezone, &c.DateFormat, &c.WeightUnit, &c.DimensionUnit, &c.Locale,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan country_config: %w", err)
		}
		items = append(items, c)
	}
	return items, nil
}
