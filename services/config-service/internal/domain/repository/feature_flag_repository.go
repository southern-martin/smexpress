package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type FeatureFlagRepository interface {
	Create(ctx context.Context, flag *entity.FeatureFlag) error
	GetByID(ctx context.Context, id string) (*entity.FeatureFlag, error)
	Update(ctx context.Context, flag *entity.FeatureFlag) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.FeatureFlag], error)
	IsEnabled(ctx context.Context, countryCode, flagKey string) (bool, error)
}
