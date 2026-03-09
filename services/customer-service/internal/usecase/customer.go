package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
	"github.com/smexpress/services/customer-service/internal/domain/repository"
)

type CustomerUseCase struct {
	repo repository.CustomerRepository
}

func NewCustomerUseCase(repo repository.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{repo: repo}
}

func (uc *CustomerUseCase) Create(ctx context.Context, customer *entity.Customer) error {
	if customer.CompanyName == "" || customer.CountryCode == "" {
		return fmt.Errorf("%w: company_name and country_code required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, customer)
}

func (uc *CustomerUseCase) GetByID(ctx context.Context, id string) (*entity.Customer, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *CustomerUseCase) Update(ctx context.Context, customer *entity.Customer) error {
	return uc.repo.Update(ctx, customer)
}

func (uc *CustomerUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *CustomerUseCase) List(ctx context.Context, countryCode, franchiseID, search string, page db.Page) (db.PagedResult[entity.Customer], error) {
	return uc.repo.List(ctx, countryCode, franchiseID, search, page)
}

func (uc *CustomerUseCase) Activate(ctx context.Context, id string) error {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	c.IsActive = true
	return uc.repo.Update(ctx, c)
}

func (uc *CustomerUseCase) Deactivate(ctx context.Context, id string) error {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	c.IsActive = false
	return uc.repo.Update(ctx, c)
}

func (uc *CustomerUseCase) SetCreditHold(ctx context.Context, id string, hold bool) error {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	c.IsCreditHold = hold
	return uc.repo.Update(ctx, c)
}
