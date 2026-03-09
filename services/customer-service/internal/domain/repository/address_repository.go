package repository

import (
	"context"

	"github.com/smexpress/services/customer-service/internal/domain/entity"
)

type AddressRepository interface {
	Create(ctx context.Context, address *entity.CustomerAddress) error
	GetByID(ctx context.Context, id string) (*entity.CustomerAddress, error)
	Update(ctx context.Context, address *entity.CustomerAddress) error
	Delete(ctx context.Context, id string) error
	ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerAddress, error)
}
