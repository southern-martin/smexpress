package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
)

type AddressRepo struct {
	pool *pgxpool.Pool
}

func NewAddressRepo(pool *pgxpool.Pool) *AddressRepo {
	return &AddressRepo{pool: pool}
}

func (r *AddressRepo) Create(ctx context.Context, a *entity.CustomerAddress) error {
	q := `INSERT INTO customer_addresses (customer_id, address_type, company_name, contact_name,
		address_line1, address_line2, city, state, postcode, country_code, phone, email, is_default, instructions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		a.CustomerID, a.AddressType, a.CompanyName, a.ContactName,
		a.AddressLine1, a.AddressLine2, a.City, a.State, a.Postcode, a.CountryCode,
		a.Phone, a.Email, a.IsDefault, a.Instructions,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AddressRepo) GetByID(ctx context.Context, id string) (*entity.CustomerAddress, error) {
	a := &entity.CustomerAddress{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, customer_id, address_type, company_name, contact_name,
		address_line1, address_line2, city, state, postcode, country_code,
		phone, email, is_default, instructions, created_at, updated_at
		FROM customer_addresses WHERE id = $1`, id,
	).Scan(&a.ID, &a.CustomerID, &a.AddressType, &a.CompanyName, &a.ContactName,
		&a.AddressLine1, &a.AddressLine2, &a.City, &a.State, &a.Postcode, &a.CountryCode,
		&a.Phone, &a.Email, &a.IsDefault, &a.Instructions, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query address: %w", err)
	}
	return a, nil
}

func (r *AddressRepo) Update(ctx context.Context, a *entity.CustomerAddress) error {
	q := `UPDATE customer_addresses SET company_name=$1, contact_name=$2, address_line1=$3, address_line2=$4,
		city=$5, state=$6, postcode=$7, country_code=$8, phone=$9, email=$10, is_default=$11, instructions=$12, updated_at=NOW()
		WHERE id=$13 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q,
		a.CompanyName, a.ContactName, a.AddressLine1, a.AddressLine2,
		a.City, a.State, a.Postcode, a.CountryCode, a.Phone, a.Email, a.IsDefault, a.Instructions, a.ID,
	).Scan(&a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update address: %w", err)
	}
	return nil
}

func (r *AddressRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM customer_addresses WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete address: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *AddressRepo) ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerAddress, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, customer_id, address_type, company_name, contact_name,
		address_line1, address_line2, city, state, postcode, country_code,
		phone, email, is_default, instructions, created_at, updated_at
		FROM customer_addresses WHERE customer_id = $1 ORDER BY is_default DESC, address_type, company_name`, customerID)
	if err != nil {
		return nil, fmt.Errorf("list addresses: %w", err)
	}
	defer rows.Close()

	var items []entity.CustomerAddress
	for rows.Next() {
		var a entity.CustomerAddress
		if err := rows.Scan(&a.ID, &a.CustomerID, &a.AddressType, &a.CompanyName, &a.ContactName,
			&a.AddressLine1, &a.AddressLine2, &a.City, &a.State, &a.Postcode, &a.CountryCode,
			&a.Phone, &a.Email, &a.IsDefault, &a.Instructions, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan address: %w", err)
		}
		items = append(items, a)
	}
	return items, nil
}
