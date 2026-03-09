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

type ZoneRepo struct {
	pool *pgxpool.Pool
}

func NewZoneRepo(pool *pgxpool.Pool) *ZoneRepo {
	return &ZoneRepo{pool: pool}
}

func (r *ZoneRepo) Create(ctx context.Context, zone *entity.Zone) error {
	q := `INSERT INTO zones (country_code, zone_name, zone_code, description)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, q, zone.CountryCode, zone.ZoneName, zone.ZoneCode, zone.Description,
	).Scan(&zone.ID, &zone.CreatedAt, &zone.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: zone code already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert zone: %w", err)
	}
	return nil
}

func (r *ZoneRepo) GetByID(ctx context.Context, id string) (*entity.Zone, error) {
	zone := &entity.Zone{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, country_code, zone_name, zone_code, description, created_at, updated_at FROM zones WHERE id = $1`, id,
	).Scan(&zone.ID, &zone.CountryCode, &zone.ZoneName, &zone.ZoneCode, &zone.Description, &zone.CreatedAt, &zone.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query zone: %w", err)
	}

	postcodes, err := r.getPostcodes(ctx, id)
	if err != nil {
		return nil, err
	}
	zone.Postcodes = postcodes
	return zone, nil
}

func (r *ZoneRepo) getPostcodes(ctx context.Context, zoneID string) ([]entity.ZonePostcode, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, zone_id, postcode_from, postcode_to, created_at FROM zone_postcodes WHERE zone_id = $1 ORDER BY postcode_from`, zoneID)
	if err != nil {
		return nil, fmt.Errorf("query zone postcodes: %w", err)
	}
	defer rows.Close()

	var items []entity.ZonePostcode
	for rows.Next() {
		var p entity.ZonePostcode
		if err := rows.Scan(&p.ID, &p.ZoneID, &p.PostcodeFrom, &p.PostcodeTo, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan zone postcode: %w", err)
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *ZoneRepo) Update(ctx context.Context, zone *entity.Zone) error {
	q := `UPDATE zones SET zone_name=$1, description=$2, updated_at=NOW() WHERE id=$3 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, zone.ZoneName, zone.Description, zone.ID).Scan(&zone.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update zone: %w", err)
	}
	return nil
}

func (r *ZoneRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM zones WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete zone: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *ZoneRepo) ListByCountry(ctx context.Context, countryCode string) ([]entity.Zone, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, country_code, zone_name, zone_code, description, created_at, updated_at FROM zones WHERE country_code = $1 ORDER BY zone_name`, countryCode)
	if err != nil {
		return nil, fmt.Errorf("list zones: %w", err)
	}
	defer rows.Close()

	var items []entity.Zone
	for rows.Next() {
		var z entity.Zone
		if err := rows.Scan(&z.ID, &z.CountryCode, &z.ZoneName, &z.ZoneCode, &z.Description, &z.CreatedAt, &z.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan zone: %w", err)
		}
		items = append(items, z)
	}
	return items, nil
}

func (r *ZoneRepo) FindZoneForPostcode(ctx context.Context, countryCode, postcode string) (*entity.Zone, error) {
	q := `SELECT z.id, z.country_code, z.zone_name, z.zone_code, z.description, z.created_at, z.updated_at
		FROM zones z
		JOIN zone_postcodes zp ON z.id = zp.zone_id
		WHERE z.country_code = $1 AND zp.postcode_from <= $2 AND zp.postcode_to >= $2
		LIMIT 1`
	zone := &entity.Zone{}
	err := r.pool.QueryRow(ctx, q, countryCode, postcode).Scan(
		&zone.ID, &zone.CountryCode, &zone.ZoneName, &zone.ZoneCode, &zone.Description, &zone.CreatedAt, &zone.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("find zone: %w", err)
	}
	return zone, nil
}

func (r *ZoneRepo) SetPostcodes(ctx context.Context, zoneID string, postcodes []entity.ZonePostcode) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM zone_postcodes WHERE zone_id = $1`, zoneID); err != nil {
		return fmt.Errorf("clear zone postcodes: %w", err)
	}

	for _, p := range postcodes {
		_, err := tx.Exec(ctx,
			`INSERT INTO zone_postcodes (zone_id, postcode_from, postcode_to) VALUES ($1, $2, $3)`,
			zoneID, p.PostcodeFrom, p.PostcodeTo)
		if err != nil {
			return fmt.Errorf("insert zone postcode: %w", err)
		}
	}

	return tx.Commit(ctx)
}
