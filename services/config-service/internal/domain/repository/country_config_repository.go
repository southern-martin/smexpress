package repository

import (
	"context"

	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type CountryConfigRepository interface {
	Create(ctx context.Context, cfg *entity.CountryConfig) error
	GetByCode(ctx context.Context, code string) (*entity.CountryConfig, error)
	Update(ctx context.Context, cfg *entity.CountryConfig) error
	List(ctx context.Context) ([]entity.CountryConfig, error)
	ListActive(ctx context.Context) ([]entity.CountryConfig, error)
}
