package repository

import (
	"context"

	"github.com/smexpress/services/customer-service/internal/domain/entity"
)

type ContactRepository interface {
	Create(ctx context.Context, contact *entity.CustomerContact) error
	GetByID(ctx context.Context, id string) (*entity.CustomerContact, error)
	Update(ctx context.Context, contact *entity.CustomerContact) error
	Delete(ctx context.Context, id string) error
	ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerContact, error)
}
