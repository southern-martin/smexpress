package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
	"github.com/smexpress/services/customer-service/internal/domain/repository"
)

type ContactUseCase struct {
	repo repository.ContactRepository
}

func NewContactUseCase(repo repository.ContactRepository) *ContactUseCase {
	return &ContactUseCase{repo: repo}
}

func (uc *ContactUseCase) Create(ctx context.Context, contact *entity.CustomerContact) error {
	if contact.CustomerID == "" || contact.FirstName == "" || contact.LastName == "" {
		return fmt.Errorf("%w: customer_id, first_name and last_name required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, contact)
}

func (uc *ContactUseCase) GetByID(ctx context.Context, id string) (*entity.CustomerContact, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ContactUseCase) Update(ctx context.Context, contact *entity.CustomerContact) error {
	return uc.repo.Update(ctx, contact)
}

func (uc *ContactUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *ContactUseCase) ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerContact, error) {
	return uc.repo.ListByCustomer(ctx, customerID)
}
