package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/user-service/internal/domain/entity"
)

type UserPreferenceRepo struct {
	pool *pgxpool.Pool
}

func NewUserPreferenceRepo(pool *pgxpool.Pool) *UserPreferenceRepo {
	return &UserPreferenceRepo{pool: pool}
}

func (r *UserPreferenceRepo) ListByUser(ctx context.Context, userID string) ([]entity.UserPreference, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, preference_key, preference_value, created_at, updated_at
		FROM user_preferences WHERE user_id = $1 ORDER BY preference_key`, userID)
	if err != nil {
		return nil, fmt.Errorf("list preferences: %w", err)
	}
	defer rows.Close()

	var items []entity.UserPreference
	for rows.Next() {
		var p entity.UserPreference
		if err := rows.Scan(&p.ID, &p.UserID, &p.PreferenceKey, &p.PreferenceValue, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan preference: %w", err)
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *UserPreferenceRepo) Set(ctx context.Context, userID, key, value string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_preferences (user_id, preference_key, preference_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, preference_key) DO UPDATE SET preference_value = $3, updated_at = NOW()`,
		userID, key, value)
	if err != nil {
		return fmt.Errorf("set preference: %w", err)
	}
	return nil
}

func (r *UserPreferenceRepo) Delete(ctx context.Context, userID, key string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM user_preferences WHERE user_id = $1 AND preference_key = $2`, userID, key)
	if err != nil {
		return fmt.Errorf("delete preference: %w", err)
	}
	return nil
}
