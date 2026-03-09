package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
)

type RefreshTokenRepo struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepo(pool *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{pool: pool}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, token *entity.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		token.UserID, token.TokenHash, token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *RefreshTokenRepo) GetByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	rt := &entity.RefreshToken{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, revoked, created_at, revoked_at
		FROM refresh_tokens WHERE token_hash = $1`, hash,
	).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.Revoked, &rt.CreatedAt, &rt.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	return rt, nil
}

func (r *RefreshTokenRepo) RevokeByUserID(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE refresh_tokens SET revoked = true, revoked_at = NOW() WHERE user_id = $1 AND revoked = false`,
		userID)
	if err != nil {
		return fmt.Errorf("revoke tokens: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE expires_at < NOW()`)
	if err != nil {
		return fmt.Errorf("delete expired tokens: %w", err)
	}
	return nil
}
