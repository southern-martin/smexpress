package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
)

type SequenceRepo struct {
	pool *pgxpool.Pool
}

func NewSequenceRepo(pool *pgxpool.Pool) *SequenceRepo {
	return &SequenceRepo{pool: pool}
}

func (r *SequenceRepo) Create(ctx context.Context, seq *entity.Sequence) error {
	query := `INSERT INTO sequences (country_code, sequence_type, prefix, format_pattern)
		VALUES ($1, $2, $3, $4)
		RETURNING id, current_value, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		seq.CountryCode, seq.SequenceType, seq.Prefix, seq.FormatPattern,
	).Scan(&seq.ID, &seq.CurrentValue, &seq.CreatedAt, &seq.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: sequence already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert sequence: %w", err)
	}
	return nil
}

func (r *SequenceRepo) List(ctx context.Context, countryCode string) ([]entity.Sequence, error) {
	query := `SELECT id, country_code, sequence_type, prefix, current_value, format_pattern, created_at, updated_at FROM sequences`
	var args []any
	if countryCode != "" {
		query += " WHERE country_code = $1"
		args = append(args, countryCode)
	}
	query += " ORDER BY country_code, sequence_type"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list sequences: %w", err)
	}
	defer rows.Close()

	var items []entity.Sequence
	for rows.Next() {
		var s entity.Sequence
		if err := rows.Scan(&s.ID, &s.CountryCode, &s.SequenceType, &s.Prefix,
			&s.CurrentValue, &s.FormatPattern, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan sequence: %w", err)
		}
		items = append(items, s)
	}
	return items, nil
}

func (r *SequenceRepo) NextValue(ctx context.Context, countryCode, seqType string) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var seq entity.Sequence
	err = tx.QueryRow(ctx,
		`SELECT id, prefix, current_value, format_pattern FROM sequences
		WHERE country_code = $1 AND sequence_type = $2 FOR UPDATE`,
		countryCode, seqType,
	).Scan(&seq.ID, &seq.Prefix, &seq.CurrentValue, &seq.FormatPattern)
	if err != nil {
		return "", fmt.Errorf("%w: sequence %s/%s", domainerr.ErrNotFound, countryCode, seqType)
	}

	newValue := seq.CurrentValue + 1
	_, err = tx.Exec(ctx,
		`UPDATE sequences SET current_value = $1, updated_at = NOW() WHERE id = $2`,
		newValue, seq.ID,
	)
	if err != nil {
		return "", fmt.Errorf("update sequence: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit tx: %w", err)
	}

	formatted := seq.FormatPattern
	formatted = strings.ReplaceAll(formatted, "{prefix}", seq.Prefix)
	formatted = strings.ReplaceAll(formatted, "{value}", fmt.Sprintf("%06d", newValue))
	return formatted, nil
}
