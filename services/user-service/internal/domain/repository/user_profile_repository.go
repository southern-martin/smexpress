package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/user-service/internal/domain/entity"
)

type UserProfileRepository interface {
	Create(ctx context.Context, profile *entity.UserProfile) error
	GetByUserID(ctx context.Context, userID string) (*entity.UserProfile, error)
	Update(ctx context.Context, profile *entity.UserProfile) error
	List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.UserProfile], error)
}

type UserPreferenceRepository interface {
	ListByUser(ctx context.Context, userID string) ([]entity.UserPreference, error)
	Set(ctx context.Context, userID, key, value string) error
	Delete(ctx context.Context, userID, key string) error
}
