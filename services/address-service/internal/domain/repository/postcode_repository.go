package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/address-service/internal/domain/entity"
)

type PostcodeRepository interface {
	Search(ctx context.Context, countryCode, query string, page db.Page) (db.PagedResult[entity.Postcode], error)
	GetByPostcode(ctx context.Context, countryCode, postcode string) ([]entity.Postcode, error)
	Create(ctx context.Context, postcode *entity.Postcode) error
	BulkCreate(ctx context.Context, postcodes []entity.Postcode) (int64, error)
}
