package repository

import (
	"context"

	"github.com/smexpress/services/customer-service/internal/domain/entity"
)

type NoteRepository interface {
	Create(ctx context.Context, note *entity.CustomerNote) error
	ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerNote, error)
	Delete(ctx context.Context, id string) error
}
