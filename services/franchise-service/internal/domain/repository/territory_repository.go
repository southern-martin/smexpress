package repository

import (
	"context"

	"github.com/smexpress/services/franchise-service/internal/domain/entity"
)

type TerritoryRepository interface {
	Create(ctx context.Context, territory *entity.Territory) error
	GetByID(ctx context.Context, id string) (*entity.Territory, error)
	Update(ctx context.Context, territory *entity.Territory) error
	Delete(ctx context.Context, id string) error
	ListByFranchise(ctx context.Context, franchiseID string) ([]entity.Territory, error)
	FindByPostcode(ctx context.Context, countryCode, postcode string) ([]entity.Territory, error)
}
