package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
)

type HolidayRepo struct {
	pool *pgxpool.Pool
}

func NewHolidayRepo(pool *pgxpool.Pool) *HolidayRepo {
	return &HolidayRepo{pool: pool}
}

func (r *HolidayRepo) Create(ctx context.Context, h *entity.Holiday) error {
	query := `INSERT INTO holidays (country_code, holiday_date, name, is_recurring)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	err := r.pool.QueryRow(ctx, query,
		h.CountryCode, h.HolidayDate, h.Name, h.IsRecurring,
	).Scan(&h.ID, &h.CreatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return fmt.Errorf("%w: holiday already exists", domainerr.ErrAlreadyExists)
		}
		return fmt.Errorf("insert holiday: %w", err)
	}
	return nil
}

func (r *HolidayRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM holidays WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete holiday: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *HolidayRepo) List(ctx context.Context, countryCode string, year int, page db.Page) (db.PagedResult[entity.Holiday], error) {
	where := "WHERE EXTRACT(YEAR FROM holiday_date) = $1"
	args := []any{year}

	if countryCode != "" {
		where += " AND country_code = $2"
		args = append(args, countryCode)
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM holidays %s", where)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return db.PagedResult[entity.Holiday]{}, fmt.Errorf("count holidays: %w", err)
	}

	limitClause, limitArgs := page.LimitOffsetClause(len(args) + 1)
	args = append(args, limitArgs...)

	dataQuery := fmt.Sprintf(`SELECT id, country_code, holiday_date, name, is_recurring, created_at
		FROM holidays %s ORDER BY holiday_date %s`, where, limitClause)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return db.PagedResult[entity.Holiday]{}, fmt.Errorf("list holidays: %w", err)
	}
	defer rows.Close()

	var items []entity.Holiday
	for rows.Next() {
		var h entity.Holiday
		if err := rows.Scan(&h.ID, &h.CountryCode, &h.HolidayDate, &h.Name, &h.IsRecurring, &h.CreatedAt); err != nil {
			return db.PagedResult[entity.Holiday]{}, fmt.Errorf("scan holiday: %w", err)
		}
		items = append(items, h)
	}

	return db.NewPagedResult(items, total, page), nil
}
