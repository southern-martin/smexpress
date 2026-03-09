package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/franchise-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
	"github.com/smexpress/services/franchise-service/internal/domain/repository"
)

type TerritoryUseCase struct {
	repo repository.TerritoryRepository
}

func NewTerritoryUseCase(repo repository.TerritoryRepository) *TerritoryUseCase {
	return &TerritoryUseCase{repo: repo}
}

func (uc *TerritoryUseCase) Create(ctx context.Context, territory *entity.Territory) error {
	if territory.FranchiseID == "" || territory.CountryCode == "" || territory.Name == "" {
		return fmt.Errorf("%w: franchise_id, country_code and name required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, territory)
}

func (uc *TerritoryUseCase) GetByID(ctx context.Context, id string) (*entity.Territory, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *TerritoryUseCase) Update(ctx context.Context, territory *entity.Territory) error {
	return uc.repo.Update(ctx, territory)
}

func (uc *TerritoryUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *TerritoryUseCase) ListByFranchise(ctx context.Context, franchiseID string) ([]entity.Territory, error) {
	return uc.repo.ListByFranchise(ctx, franchiseID)
}

func (uc *TerritoryUseCase) FindByPostcode(ctx context.Context, countryCode, postcode string) ([]entity.Territory, error) {
	if countryCode == "" || postcode == "" {
		return nil, fmt.Errorf("%w: country_code and postcode required", domainerr.ErrInvalidInput)
	}
	return uc.repo.FindByPostcode(ctx, countryCode, postcode)
}
