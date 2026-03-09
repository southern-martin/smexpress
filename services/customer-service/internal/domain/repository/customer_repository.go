package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/customer-service/internal/domain/entity"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	GetByID(ctx context.Context, id string) (*entity.Customer, error)
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, countryCode, franchiseID, search string, page db.Page) (db.PagedResult[entity.Customer], error)
}
