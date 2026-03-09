package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/address-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/address-service/internal/domain/errors"
)

type RegionRepo struct {
	pool *pgxpool.Pool
}

func NewRegionRepo(pool *pgxpool.Pool) *RegionRepo {
	return &RegionRepo{pool: pool}
}

func (r *RegionRepo) Create(ctx context.Context, region *entity.Region) error {
	q := `INSERT INTO regions (country_code, name, code, parent_region_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, q, region.CountryCode, region.Name, region.Code, region.ParentRegionID,
	).Scan(&region.ID, &region.CreatedAt, &region.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: region code already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert region: %w", err)
	}
	return nil
}

func (r *RegionRepo) GetByID(ctx context.Context, id string) (*entity.Region, error) {
	region := &entity.Region{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, country_code, name, code, parent_region_id, created_at, updated_at FROM regions WHERE id = $1`, id,
	).Scan(&region.ID, &region.CountryCode, &region.Name, &region.Code, &region.ParentRegionID, &region.CreatedAt, &region.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query region: %w", err)
	}
	return region, nil
}

func (r *RegionRepo) Update(ctx context.Context, region *entity.Region) error {
	q := `UPDATE regions SET name=$1, updated_at=NOW() WHERE id=$2 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, region.Name, region.ID).Scan(&region.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update region: %w", err)
	}
	return nil
}

func (r *RegionRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM regions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete region: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *RegionRepo) ListByCountry(ctx context.Context, countryCode string) ([]entity.Region, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, country_code, name, code, parent_region_id, created_at, updated_at FROM regions WHERE country_code = $1 ORDER BY name`, countryCode)
	if err != nil {
		return nil, fmt.Errorf("list regions: %w", err)
	}
	defer rows.Close()

	var items []entity.Region
	for rows.Next() {
		var reg entity.Region
		if err := rows.Scan(&reg.ID, &reg.CountryCode, &reg.Name, &reg.Code, &reg.ParentRegionID, &reg.CreatedAt, &reg.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan region: %w", err)
		}
		items = append(items, reg)
	}
	return items, nil
}
