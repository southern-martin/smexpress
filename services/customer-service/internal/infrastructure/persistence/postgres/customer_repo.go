package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
)

type CustomerRepo struct {
	pool *pgxpool.Pool
}

func NewCustomerRepo(pool *pgxpool.Pool) *CustomerRepo {
	return &CustomerRepo{pool: pool}
}

func (r *CustomerRepo) Create(ctx context.Context, c *entity.Customer) error {
	q := `INSERT INTO customers (country_code, franchise_id, company_name, trading_name, account_number,
		abn, email, phone, website, credit_limit, payment_terms, created_by)
		VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7, $8, $9, $10, $11, NULLIF($12, ''))
		RETURNING id, credit_balance, is_active, is_credit_hold, created_at, updated_at`
	err := r.pool.QueryRow(ctx, q,
		c.CountryCode, c.FranchiseID, c.CompanyName, c.TradingName, c.AccountNumber,
		c.ABN, c.Email, c.Phone, c.Website, c.CreditLimit, c.PaymentTerms, c.CreatedBy,
	).Scan(&c.ID, &c.CreditBalance, &c.IsActive, &c.IsCreditHold, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: account number already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert customer: %w", err)
	}
	return nil
}

func (r *CustomerRepo) GetByID(ctx context.Context, id string) (*entity.Customer, error) {
	q := `SELECT id, country_code, COALESCE(franchise_id::text, ''), company_name, trading_name,
		COALESCE(account_number, ''), abn, email, phone, website,
		credit_limit, credit_balance, payment_terms, is_active, is_credit_hold,
		COALESCE(notes, ''), COALESCE(created_by::text, ''), created_at, updated_at
		FROM customers WHERE id = $1`
	c := &entity.Customer{}
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&c.ID, &c.CountryCode, &c.FranchiseID, &c.CompanyName, &c.TradingName,
		&c.AccountNumber, &c.ABN, &c.Email, &c.Phone, &c.Website,
		&c.CreditLimit, &c.CreditBalance, &c.PaymentTerms, &c.IsActive, &c.IsCreditHold,
		&c.Notes, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("query customer: %w", err)
	}
	return c, nil
}

func (r *CustomerRepo) Update(ctx context.Context, c *entity.Customer) error {
	q := `UPDATE customers SET company_name=$1, trading_name=$2, abn=$3, email=$4, phone=$5,
		website=$6, credit_limit=$7, payment_terms=$8, is_active=$9, is_credit_hold=$10, updated_at=NOW()
		WHERE id=$11 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q,
		c.CompanyName, c.TradingName, c.ABN, c.Email, c.Phone,
		c.Website, c.CreditLimit, c.PaymentTerms, c.IsActive, c.IsCreditHold, c.ID,
	).Scan(&c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return fmt.Errorf("update customer: %w", err)
	}
	return nil
}

func (r *CustomerRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM customers WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *CustomerRepo) List(ctx context.Context, countryCode, franchiseID, search string, page db.Page) (db.PagedResult[entity.Customer], error) {
	where := "WHERE 1=1"
	args := []any{}
	argIdx := 0

	if countryCode != "" {
		argIdx++
		where += fmt.Sprintf(" AND country_code = $%d", argIdx)
		args = append(args, countryCode)
	}
	if franchiseID != "" {
		argIdx++
		where += fmt.Sprintf(" AND franchise_id = $%d", argIdx)
		args = append(args, franchiseID)
	}
	if search != "" {
		argIdx++
		where += fmt.Sprintf(" AND (company_name ILIKE $%d OR trading_name ILIKE $%d OR account_number ILIKE $%d OR email ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM customers "+where, args...).Scan(&total); err != nil {
		return db.PagedResult[entity.Customer]{}, fmt.Errorf("count customers: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(argIdx + 1)
	args = append(args, limitArgs...)

	q := fmt.Sprintf(`SELECT id, country_code, COALESCE(franchise_id::text, ''), company_name, trading_name,
		COALESCE(account_number, ''), abn, email, phone, website,
		credit_limit, credit_balance, payment_terms, is_active, is_credit_hold,
		COALESCE(notes, ''), COALESCE(created_by::text, ''), created_at, updated_at
		FROM customers %s ORDER BY company_name %s`, where, limitClause)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return db.PagedResult[entity.Customer]{}, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	var items []entity.Customer
	for rows.Next() {
		var c entity.Customer
		if err := rows.Scan(
			&c.ID, &c.CountryCode, &c.FranchiseID, &c.CompanyName, &c.TradingName,
			&c.AccountNumber, &c.ABN, &c.Email, &c.Phone, &c.Website,
			&c.CreditLimit, &c.CreditBalance, &c.PaymentTerms, &c.IsActive, &c.IsCreditHold,
			&c.Notes, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return db.PagedResult[entity.Customer]{}, fmt.Errorf("scan customer: %w", err)
		}
		items = append(items, c)
	}

	return db.NewPagedResult(items, total, page), nil
}
