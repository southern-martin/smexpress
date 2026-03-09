package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type SystemConfigRepository interface {
	Create(ctx context.Context, cfg *entity.SystemConfig) error
	GetByID(ctx context.Context, id string) (*entity.SystemConfig, error)
	GetByKey(ctx context.Context, countryCode, key string) (*entity.SystemConfig, error)
	Update(ctx context.Context, cfg *entity.SystemConfig) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.SystemConfig], error)
}
