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

type ContactRepo struct {
	pool *pgxpool.Pool
}

func NewContactRepo(pool *pgxpool.Pool) *ContactRepo {
	return &ContactRepo{pool: pool}
}

func (r *ContactRepo) Create(ctx context.Context, c *entity.CustomerContact) error {
	q := `INSERT INTO customer_contacts (customer_id, first_name, last_name, email, phone, mobile, position, is_primary, is_billing)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		c.CustomerID, c.FirstName, c.LastName, c.Email, c.Phone, c.Mobile, c.Position, c.IsPrimary, c.IsBilling,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *ContactRepo) GetByID(ctx context.Context, id string) (*entity.CustomerContact, error) {
	c := &entity.CustomerContact{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, customer_id, first_name, last_name, email, phone, mobile, position, is_primary, is_billing, created_at, updated_at
		FROM customer_contacts WHERE id = $1`, id,
	).Scan(&c.ID, &c.CustomerID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Mobile, &c.Position, &c.IsPrimary, &c.IsBilling, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query contact: %w", err)
	}
	return c, nil
}

func (r *ContactRepo) Update(ctx context.Context, c *entity.CustomerContact) error {
	q := `UPDATE customer_contacts SET first_name=$1, last_name=$2, email=$3, phone=$4, mobile=$5,
		position=$6, is_primary=$7, is_billing=$8, updated_at=NOW() WHERE id=$9 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q,
		c.FirstName, c.LastName, c.Email, c.Phone, c.Mobile, c.Position, c.IsPrimary, c.IsBilling, c.ID,
	).Scan(&c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update contact: %w", err)
	}
	return nil
}

func (r *ContactRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM customer_contacts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete contact: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *ContactRepo) ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerContact, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, customer_id, first_name, last_name, email, phone, mobile, position, is_primary, is_billing, created_at, updated_at
		FROM customer_contacts WHERE customer_id = $1 ORDER BY is_primary DESC, first_name`, customerID)
	if err != nil {
		return nil, fmt.Errorf("list contacts: %w", err)
	}
	defer rows.Close()

	var items []entity.CustomerContact
	for rows.Next() {
		var c entity.CustomerContact
		if err := rows.Scan(&c.ID, &c.CustomerID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Mobile,
			&c.Position, &c.IsPrimary, &c.IsBilling, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan contact: %w", err)
		}
		items = append(items, c)
	}
	return items, nil
}
