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

type LedgerRepo struct {
	pool *pgxpool.Pool
}

func NewLedgerRepo(pool *pgxpool.Pool) *LedgerRepo {
	return &LedgerRepo{pool: pool}
}

func (r *LedgerRepo) GetByFranchiseID(ctx context.Context, franchiseID string) (*entity.FranchiseLedger, error) {
	query := `SELECT id, franchise_id, country_code, currency, balance, created_at, updated_at
		FROM franchise_ledgers WHERE franchise_id = $1`
	l := &entity.FranchiseLedger{}
	err := r.pool.QueryRow(ctx, query, franchiseID).Scan(
		&l.ID, &l.FranchiseID, &l.CountryCode, &l.Currency, &l.Balance, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query ledger: %w", err)
	}
	return l, nil
}

func (r *LedgerRepo) Create(ctx context.Context, l *entity.FranchiseLedger) error {
	query := `INSERT INTO franchise_ledgers (franchise_id, country_code, currency, balance)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query,
		l.FranchiseID, l.CountryCode, l.Currency, l.Balance,
	).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: ledger already exists for this franchise", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert ledger: %w", err)
	}
	return nil
}

func (r *LedgerRepo) AddEntry(ctx context.Context, ledgerID string, entry *entity.FranchiseLedgerEntry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Lock the ledger row and update balance
	var currentBalance float64
	err = tx.QueryRow(ctx, `SELECT balance FROM franchise_ledgers WHERE id = $1 FOR UPDATE`, ledgerID).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("lock ledger: %w", err)
	}

	newBalance := currentBalance + entry.Amount
	if entry.EntryType == "debit" && newBalance < 0 {
		return fmt.Errorf("%w: insufficient balance", domainerr.ErrInsufficientBalance)
	}

	entry.BalanceAfter = newBalance

	insertQuery := `INSERT INTO franchise_ledger_entries (ledger_id, entry_type, amount, balance_after, description, reference_type, reference_id)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''))
		RETURNING id, created_at`
	err = tx.QueryRow(ctx, insertQuery,
		ledgerID, entry.EntryType, entry.Amount, entry.BalanceAfter,
		entry.Description, entry.ReferenceType, entry.ReferenceID,
	).Scan(&entry.ID, &entry.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert ledger entry: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE franchise_ledgers SET balance = $1, updated_at = NOW() WHERE id = $2`, newBalance, ledgerID)
	if err != nil {
		return fmt.Errorf("update ledger balance: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *LedgerRepo) ListEntries(ctx context.Context, ledgerID string, page db.Page) (db.PagedResult[entity.FranchiseLedgerEntry], error) {
	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM franchise_ledger_entries WHERE ledger_id = $1", ledgerID).Scan(&total); err != nil {
		return db.PagedResult[entity.FranchiseLedgerEntry]{}, fmt.Errorf("count ledger entries: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(2)
	args := []any{ledgerID}
	args = append(args, limitArgs...)

	query := fmt.Sprintf(`SELECT id, ledger_id, entry_type, amount, balance_after, description,
		COALESCE(reference_type, ''), COALESCE(reference_id::text, ''), created_at
		FROM franchise_ledger_entries WHERE ledger_id = $1 ORDER BY created_at DESC %s`, limitClause)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return db.PagedResult[entity.FranchiseLedgerEntry]{}, fmt.Errorf("list ledger entries: %w", err)
	}
	defer rows.Close()

	var items []entity.FranchiseLedgerEntry
	for rows.Next() {
		var e entity.FranchiseLedgerEntry
		if err := rows.Scan(&e.ID, &e.LedgerID, &e.EntryType, &e.Amount, &e.BalanceAfter,
			&e.Description, &e.ReferenceType, &e.ReferenceID, &e.CreatedAt); err != nil {
			return db.PagedResult[entity.FranchiseLedgerEntry]{}, fmt.Errorf("scan ledger entry: %w", err)
		}
		items = append(items, e)
	}

	return db.NewPagedResult(items, total, page), nil
}
