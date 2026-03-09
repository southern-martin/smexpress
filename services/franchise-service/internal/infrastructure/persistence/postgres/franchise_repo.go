package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
)

type FranchiseRepo struct {
	pool *pgxpool.Pool
}

func NewFranchiseRepo(pool *pgxpool.Pool) *FranchiseRepo {
	return &FranchiseRepo{pool: pool}
}

func (r *FranchiseRepo) Create(ctx context.Context, f *entity.Franchise) error {
	query := `INSERT INTO franchises (country_code, name, code, contact_name, email, phone,
		address_line1, address_line2, city, state, postcode, commission_rate, parent_franchise_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, is_active, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		f.CountryCode, f.Name, f.Code, f.ContactName, f.Email, f.Phone,
		f.AddressLine1, f.AddressLine2, f.City, f.State, f.Postcode,
		f.CommissionRate, f.ParentFranchiseID,
	).Scan(&f.ID, &f.IsActive, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: franchise code already exists in this country", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert franchise: %w", err)
	}
	return nil
}

func (r *FranchiseRepo) GetByID(ctx context.Context, id string) (*entity.Franchise, error) {
	query := `SELECT id, country_code, name, code, contact_name, email, phone,
		address_line1, address_line2, city, state, postcode, is_active, commission_rate,
		parent_franchise_id, created_at, updated_at
		FROM franchises WHERE id = $1`
	return r.scanFranchise(ctx, query, id)
}

func (r *FranchiseRepo) GetByCode(ctx context.Context, countryCode, code string) (*entity.Franchise, error) {
	query := `SELECT id, country_code, name, code, contact_name, email, phone,
		address_line1, address_line2, city, state, postcode, is_active, commission_rate,
		parent_franchise_id, created_at, updated_at
		FROM franchises WHERE country_code = $1 AND code = $2`
	return r.scanFranchise(ctx, query, countryCode, code)
}

func (r *FranchiseRepo) scanFranchise(ctx context.Context, query string, args ...any) (*entity.Franchise, error) {
	f := &entity.Franchise{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&f.ID, &f.CountryCode, &f.Name, &f.Code, &f.ContactName, &f.Email, &f.Phone,
		&f.AddressLine1, &f.AddressLine2, &f.City, &f.State, &f.Postcode,
		&f.IsActive, &f.CommissionRate, &f.ParentFranchiseID, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query franchise: %w", err)
	}
	return f, nil
}

func (r *FranchiseRepo) Update(ctx context.Context, f *entity.Franchise) error {
	query := `UPDATE franchises SET name=$1, contact_name=$2, email=$3, phone=$4,
		address_line1=$5, address_line2=$6, city=$7, state=$8, postcode=$9,
		is_active=$10, commission_rate=$11, updated_at=NOW()
		WHERE id = $12 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query,
		f.Name, f.ContactName, f.Email, f.Phone,
		f.AddressLine1, f.AddressLine2, f.City, f.State, f.Postcode,
		f.IsActive, f.CommissionRate, f.ID,
	).Scan(&f.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update franchise: %w", err)
	}
	return nil
}

func (r *FranchiseRepo) List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.Franchise], error) {
	where := "WHERE 1=1"
	args := []any{}
	argIdx := 0

	if countryCode != "" {
		argIdx++
		where += fmt.Sprintf(" AND country_code = $%d", argIdx)
		args = append(args, countryCode)
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM franchises "+where, args...).Scan(&total); err != nil {
		return db.PagedResult[entity.Franchise]{}, fmt.Errorf("count franchises: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(argIdx + 1)
	args = append(args, limitArgs...)

	query := fmt.Sprintf(`SELECT id, country_code, name, code, contact_name, email, phone,
		address_line1, address_line2, city, state, postcode, is_active, commission_rate,
		parent_franchise_id, created_at, updated_at
		FROM franchises %s ORDER BY name %s`, where, limitClause)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return db.PagedResult[entity.Franchise]{}, fmt.Errorf("list franchises: %w", err)
	}
	defer rows.Close()

	var items []entity.Franchise
	for rows.Next() {
		var f entity.Franchise
		if err := rows.Scan(
			&f.ID, &f.CountryCode, &f.Name, &f.Code, &f.ContactName, &f.Email, &f.Phone,
			&f.AddressLine1, &f.AddressLine2, &f.City, &f.State, &f.Postcode,
			&f.IsActive, &f.CommissionRate, &f.ParentFranchiseID, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return db.PagedResult[entity.Franchise]{}, fmt.Errorf("scan franchise: %w", err)
		}
		items = append(items, f)
	}

	return db.NewPagedResult(items, total, page), nil
}

func (r *FranchiseRepo) ListByCountry(ctx context.Context, countryCode string) ([]entity.Franchise, error) {
	query := `SELECT id, country_code, name, code, contact_name, email, phone,
		address_line1, address_line2, city, state, postcode, is_active, commission_rate,
		parent_franchise_id, created_at, updated_at
		FROM franchises WHERE country_code = $1 ORDER BY name`

	rows, err := r.pool.Query(ctx, query, countryCode)
	if err != nil {
		return nil, fmt.Errorf("list franchises by country: %w", err)
	}
	defer rows.Close()

	var items []entity.Franchise
	for rows.Next() {
		var f entity.Franchise
		if err := rows.Scan(
			&f.ID, &f.CountryCode, &f.Name, &f.Code, &f.ContactName, &f.Email, &f.Phone,
			&f.AddressLine1, &f.AddressLine2, &f.City, &f.State, &f.Postcode,
			&f.IsActive, &f.CommissionRate, &f.ParentFranchiseID, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan franchise: %w", err)
		}
		items = append(items, f)
	}
	return items, nil
}

func (r *FranchiseRepo) GetSettings(ctx context.Context, franchiseID string) ([]entity.FranchiseSetting, error) {
	query := `SELECT id, franchise_id, setting_key, setting_value, created_at, updated_at
		FROM franchise_settings WHERE franchise_id = $1 ORDER BY setting_key`
	rows, err := r.pool.Query(ctx, query, franchiseID)
	if err != nil {
		return nil, fmt.Errorf("list franchise settings: %w", err)
	}
	defer rows.Close()

	var items []entity.FranchiseSetting
	for rows.Next() {
		var s entity.FranchiseSetting
		if err := rows.Scan(&s.ID, &s.FranchiseID, &s.SettingKey, &s.SettingValue, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan franchise setting: %w", err)
		}
		items = append(items, s)
	}
	return items, nil
}

func (r *FranchiseRepo) SetSetting(ctx context.Context, franchiseID, key, value string) error {
	query := `INSERT INTO franchise_settings (franchise_id, setting_key, setting_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (franchise_id, setting_key)
		DO UPDATE SET setting_value = EXCLUDED.setting_value, updated_at = NOW()`
	_, err := r.pool.Exec(ctx, query, franchiseID, key, value)
	if err != nil {
		return fmt.Errorf("set franchise setting: %w", err)
	}
	return nil
}
