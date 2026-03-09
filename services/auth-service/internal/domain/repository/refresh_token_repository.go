package repository

import (
	"context"

	"github.com/smexpress/services/auth-service/internal/domain/entity"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*entity.RefreshToken, error)
	RevokeByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}
