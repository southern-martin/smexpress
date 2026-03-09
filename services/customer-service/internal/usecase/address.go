package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
	"github.com/smexpress/services/customer-service/internal/domain/repository"
)

type AddressUseCase struct {
	repo repository.AddressRepository
}

func NewAddressUseCase(repo repository.AddressRepository) *AddressUseCase {
	return &AddressUseCase{repo: repo}
}

func (uc *AddressUseCase) Create(ctx context.Context, addr *entity.CustomerAddress) error {
	if addr.CustomerID == "" || addr.AddressLine1 == "" || addr.City == "" || addr.Postcode == "" || addr.CountryCode == "" {
		return fmt.Errorf("%w: customer_id, address_line1, city, postcode and country_code required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, addr)
}

func (uc *AddressUseCase) GetByID(ctx context.Context, id string) (*entity.CustomerAddress, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *AddressUseCase) Update(ctx context.Context, addr *entity.CustomerAddress) error {
	return uc.repo.Update(ctx, addr)
}

func (uc *AddressUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *AddressUseCase) ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerAddress, error) {
	return uc.repo.ListByCustomer(ctx, customerID)
}
