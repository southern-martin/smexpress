package repository

import (
	"context"

	"github.com/smexpress/services/address-service/internal/domain/entity"
)

type RegionRepository interface {
	Create(ctx context.Context, region *entity.Region) error
	GetByID(ctx context.Context, id string) (*entity.Region, error)
	Update(ctx context.Context, region *entity.Region) error
	Delete(ctx context.Context, id string) error
	ListByCountry(ctx context.Context, countryCode string) ([]entity.Region, error)
}
