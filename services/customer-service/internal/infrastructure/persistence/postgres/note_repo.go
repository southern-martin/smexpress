package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
)

type NoteRepo struct {
	pool *pgxpool.Pool
}

func NewNoteRepo(pool *pgxpool.Pool) *NoteRepo {
	return &NoteRepo{pool: pool}
}

func (r *NoteRepo) Create(ctx context.Context, n *entity.CustomerNote) error {
	q := `INSERT INTO customer_notes (customer_id, note, created_by) VALUES ($1, $2, NULLIF($3, '')) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, q, n.CustomerID, n.Note, n.CreatedBy).Scan(&n.ID, &n.CreatedAt)
}

func (r *NoteRepo) ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerNote, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, customer_id, note, COALESCE(created_by::text, ''), created_at FROM customer_notes WHERE customer_id = $1 ORDER BY created_at DESC`, customerID)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	defer rows.Close()

	var items []entity.CustomerNote
	for rows.Next() {
		var n entity.CustomerNote
		if err := rows.Scan(&n.ID, &n.CustomerID, &n.Note, &n.CreatedBy, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		items = append(items, n)
	}
	return items, nil
}

func (r *NoteRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM customer_notes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}
