package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/address-service/internal/domain/entity"
)

type PostcodeRepo struct {
	pool *pgxpool.Pool
}

func NewPostcodeRepo(pool *pgxpool.Pool) *PostcodeRepo {
	return &PostcodeRepo{pool: pool}
}

func (r *PostcodeRepo) Search(ctx context.Context, countryCode, query string, page db.Page) (db.PagedResult[entity.Postcode], error) {
	where := "WHERE country_code = $1 AND (postcode ILIKE $2 OR suburb ILIKE $2 OR city ILIKE $2)"
	args := []any{countryCode, "%" + query + "%"}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM postcodes "+where, args...).Scan(&total); err != nil {
		return db.PagedResult[entity.Postcode]{}, fmt.Errorf("count postcodes: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(3)
	args = append(args, limitArgs...)

	q := fmt.Sprintf(`SELECT id, country_code, postcode, suburb, city, state, state_code, latitude, longitude, created_at, updated_at
		FROM postcodes %s ORDER BY postcode, suburb %s`, where, limitClause)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return db.PagedResult[entity.Postcode]{}, fmt.Errorf("search postcodes: %w", err)
	}
	defer rows.Close()

	var items []entity.Postcode
	for rows.Next() {
		var p entity.Postcode
		if err := rows.Scan(&p.ID, &p.CountryCode, &p.Postcode, &p.Suburb, &p.City, &p.State, &p.StateCode,
			&p.Latitude, &p.Longitude, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return db.PagedResult[entity.Postcode]{}, fmt.Errorf("scan postcode: %w", err)
		}
		items = append(items, p)
	}
	return db.NewPagedResult(items, total, page), nil
}

func (r *PostcodeRepo) GetByPostcode(ctx context.Context, countryCode, postcode string) ([]entity.Postcode, error) {
	q := `SELECT id, country_code, postcode, suburb, city, state, state_code, latitude, longitude, created_at, updated_at
		FROM postcodes WHERE country_code = $1 AND postcode = $2 ORDER BY suburb`

	rows, err := r.pool.Query(ctx, q, countryCode, postcode)
	if err != nil {
		return nil, fmt.Errorf("lookup postcode: %w", err)
	}
	defer rows.Close()

	var items []entity.Postcode
	for rows.Next() {
		var p entity.Postcode
		if err := rows.Scan(&p.ID, &p.CountryCode, &p.Postcode, &p.Suburb, &p.City, &p.State, &p.StateCode,
			&p.Latitude, &p.Longitude, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan postcode: %w", err)
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *PostcodeRepo) Create(ctx context.Context, p *entity.Postcode) error {
	q := `INSERT INTO postcodes (country_code, postcode, suburb, city, state, state_code, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		p.CountryCode, p.Postcode, p.Suburb, p.City, p.State, p.StateCode, p.Latitude, p.Longitude,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PostcodeRepo) BulkCreate(ctx context.Context, postcodes []entity.Postcode) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var count int64
	for _, p := range postcodes {
		_, err := tx.Exec(ctx,
			`INSERT INTO postcodes (country_code, postcode, suburb, city, state, state_code, latitude, longitude)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING`,
			p.CountryCode, p.Postcode, p.Suburb, p.City, p.State, p.StateCode, p.Latitude, p.Longitude,
		)
		if err != nil {
			return 0, fmt.Errorf("insert postcode: %w", err)
		}
		count++
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return count, nil
}
