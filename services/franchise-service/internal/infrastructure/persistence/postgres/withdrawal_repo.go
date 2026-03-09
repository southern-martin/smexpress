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

type WithdrawalRepo struct {
	pool *pgxpool.Pool
}

func NewWithdrawalRepo(pool *pgxpool.Pool) *WithdrawalRepo {
	return &WithdrawalRepo{pool: pool}
}

func (r *WithdrawalRepo) Create(ctx context.Context, w *entity.FranchiseWithdrawal) error {
	query := `INSERT INTO franchise_withdrawals (franchise_id, country_code, amount, currency, status, requested_by,
		bank_account_name, bank_account_number, bank_bsb, notes, requested_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`
	err := r.pool.QueryRow(ctx, query,
		w.FranchiseID, w.CountryCode, w.Amount, w.Currency, w.Status, w.RequestedBy,
		w.BankAccountName, w.BankAccountNumber, w.BankBSB, w.Notes, w.RequestedAt,
	).Scan(&w.ID)
	if err != nil {
		return fmt.Errorf("insert withdrawal: %w", err)
	}
	return nil
}

func (r *WithdrawalRepo) GetByID(ctx context.Context, id string) (*entity.FranchiseWithdrawal, error) {
	query := `SELECT id, franchise_id, country_code, amount, currency, status, requested_by,
		approved_by, bank_account_name, bank_account_number, bank_bsb, notes, requested_at, processed_at
		FROM franchise_withdrawals WHERE id = $1`
	w := &entity.FranchiseWithdrawal{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&w.ID, &w.FranchiseID, &w.CountryCode, &w.Amount, &w.Currency, &w.Status, &w.RequestedBy,
		&w.ApprovedBy, &w.BankAccountName, &w.BankAccountNumber, &w.BankBSB, &w.Notes,
		&w.RequestedAt, &w.ProcessedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query withdrawal: %w", err)
	}
	return w, nil
}

func (r *WithdrawalRepo) Update(ctx context.Context, w *entity.FranchiseWithdrawal) error {
	query := `UPDATE franchise_withdrawals SET status=$1, approved_by=$2, processed_at=$3
		WHERE id = $4`
	tag, err := r.pool.Exec(ctx, query, w.Status, w.ApprovedBy, w.ProcessedAt, w.ID)
	if err != nil {
		return fmt.Errorf("update withdrawal: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *WithdrawalRepo) ListByFranchise(ctx context.Context, franchiseID string, page db.Page) (db.PagedResult[entity.FranchiseWithdrawal], error) {
	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM franchise_withdrawals WHERE franchise_id = $1", franchiseID).Scan(&total); err != nil {
		return db.PagedResult[entity.FranchiseWithdrawal]{}, fmt.Errorf("count withdrawals: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(2)
	args := []any{franchiseID}
	args = append(args, limitArgs...)

	query := fmt.Sprintf(`SELECT id, franchise_id, country_code, amount, currency, status, requested_by,
		approved_by, bank_account_name, bank_account_number, bank_bsb, notes, requested_at, processed_at
		FROM franchise_withdrawals WHERE franchise_id = $1 ORDER BY requested_at DESC %s`, limitClause)

	return r.queryWithdrawals(ctx, query, args, total, page)
}

func (r *WithdrawalRepo) ListPending(ctx context.Context, page db.Page) (db.PagedResult[entity.FranchiseWithdrawal], error) {
	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM franchise_withdrawals WHERE status = 'pending'").Scan(&total); err != nil {
		return db.PagedResult[entity.FranchiseWithdrawal]{}, fmt.Errorf("count pending withdrawals: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(1)
	args := limitArgs

	query := fmt.Sprintf(`SELECT id, franchise_id, country_code, amount, currency, status, requested_by,
		approved_by, bank_account_name, bank_account_number, bank_bsb, notes, requested_at, processed_at
		FROM franchise_withdrawals WHERE status = 'pending' ORDER BY requested_at ASC %s`, limitClause)

	return r.queryWithdrawals(ctx, query, args, total, page)
}

func (r *WithdrawalRepo) queryWithdrawals(ctx context.Context, query string, args []any, total int64, page db.Page) (db.PagedResult[entity.FranchiseWithdrawal], error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return db.PagedResult[entity.FranchiseWithdrawal]{}, fmt.Errorf("list withdrawals: %w", err)
	}
	defer rows.Close()

	var items []entity.FranchiseWithdrawal
	for rows.Next() {
		var w entity.FranchiseWithdrawal
		if err := rows.Scan(
			&w.ID, &w.FranchiseID, &w.CountryCode, &w.Amount, &w.Currency, &w.Status, &w.RequestedBy,
			&w.ApprovedBy, &w.BankAccountName, &w.BankAccountNumber, &w.BankBSB, &w.Notes,
			&w.RequestedAt, &w.ProcessedAt,
		); err != nil {
			return db.PagedResult[entity.FranchiseWithdrawal]{}, fmt.Errorf("scan withdrawal: %w", err)
		}
		items = append(items, w)
	}

	return db.NewPagedResult(items, total, page), nil
}
