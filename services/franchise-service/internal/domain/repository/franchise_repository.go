package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
)

type FranchiseRepository interface {
	Create(ctx context.Context, franchise *entity.Franchise) error
	GetByID(ctx context.Context, id string) (*entity.Franchise, error)
	GetByCode(ctx context.Context, countryCode, code string) (*entity.Franchise, error)
	Update(ctx context.Context, franchise *entity.Franchise) error
	List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.Franchise], error)
	ListByCountry(ctx context.Context, countryCode string) ([]entity.Franchise, error)
	GetSettings(ctx context.Context, franchiseID string) ([]entity.FranchiseSetting, error)
	SetSetting(ctx context.Context, franchiseID, key, value string) error
}
