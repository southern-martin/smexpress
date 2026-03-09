package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
)

type TerritoryRepo struct {
	pool *pgxpool.Pool
}

func NewTerritoryRepo(pool *pgxpool.Pool) *TerritoryRepo {
	return &TerritoryRepo{pool: pool}
}

func (r *TerritoryRepo) Create(ctx context.Context, t *entity.Territory) error {
	query := `INSERT INTO territories (franchise_id, country_code, name, postcode_from, postcode_to, state, is_exclusive)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		t.FranchiseID, t.CountryCode, t.Name, t.PostcodeFrom, t.PostcodeTo, t.State, t.IsExclusive,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert territory: %w", err)
	}
	return nil
}

func (r *TerritoryRepo) GetByID(ctx context.Context, id string) (*entity.Territory, error) {
	query := `SELECT id, franchise_id, country_code, name, postcode_from, postcode_to, state, is_exclusive, created_at, updated_at
		FROM territories WHERE id = $1`
	t := &entity.Territory{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.FranchiseID, &t.CountryCode, &t.Name, &t.PostcodeFrom, &t.PostcodeTo,
		&t.State, &t.IsExclusive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query territory: %w", err)
	}
	return t, nil
}

func (r *TerritoryRepo) Update(ctx context.Context, t *entity.Territory) error {
	query := `UPDATE territories SET name=$1, postcode_from=$2, postcode_to=$3, state=$4, is_exclusive=$5, updated_at=NOW()
		WHERE id = $6 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query,
		t.Name, t.PostcodeFrom, t.PostcodeTo, t.State, t.IsExclusive, t.ID,
	).Scan(&t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update territory: %w", err)
	}
	return nil
}

func (r *TerritoryRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM territories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete territory: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *TerritoryRepo) ListByFranchise(ctx context.Context, franchiseID string) ([]entity.Territory, error) {
	query := `SELECT id, franchise_id, country_code, name, postcode_from, postcode_to, state, is_exclusive, created_at, updated_at
		FROM territories WHERE franchise_id = $1 ORDER BY name`
	rows, err := r.pool.Query(ctx, query, franchiseID)
	if err != nil {
		return nil, fmt.Errorf("list territories: %w", err)
	}
	defer rows.Close()

	var items []entity.Territory
	for rows.Next() {
		var t entity.Territory
		if err := rows.Scan(&t.ID, &t.FranchiseID, &t.CountryCode, &t.Name, &t.PostcodeFrom, &t.PostcodeTo,
			&t.State, &t.IsExclusive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan territory: %w", err)
		}
		items = append(items, t)
	}
	return items, nil
}

func (r *TerritoryRepo) FindByPostcode(ctx context.Context, countryCode, postcode string) ([]entity.Territory, error) {
	query := `SELECT id, franchise_id, country_code, name, postcode_from, postcode_to, state, is_exclusive, created_at, updated_at
		FROM territories WHERE country_code = $1 AND postcode_from <= $2 AND postcode_to >= $2 ORDER BY name`
	rows, err := r.pool.Query(ctx, query, countryCode, postcode)
	if err != nil {
		return nil, fmt.Errorf("find territories by postcode: %w", err)
	}
	defer rows.Close()

	var items []entity.Territory
	for rows.Next() {
		var t entity.Territory
		if err := rows.Scan(&t.ID, &t.FranchiseID, &t.CountryCode, &t.Name, &t.PostcodeFrom, &t.PostcodeTo,
			&t.State, &t.IsExclusive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan territory: %w", err)
		}
		items = append(items, t)
	}
	return items, nil
}
