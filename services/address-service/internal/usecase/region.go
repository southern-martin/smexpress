package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/address-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/address-service/internal/domain/errors"
	"github.com/smexpress/services/address-service/internal/domain/repository"
)

type RegionUseCase struct {
	repo repository.RegionRepository
}

func NewRegionUseCase(repo repository.RegionRepository) *RegionUseCase {
	return &RegionUseCase{repo: repo}
}

func (uc *RegionUseCase) Create(ctx context.Context, region *entity.Region) error {
	if region.CountryCode == "" || region.Name == "" || region.Code == "" {
		return fmt.Errorf("%w: country_code, name and code required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, region)
}

func (uc *RegionUseCase) GetByID(ctx context.Context, id string) (*entity.Region, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *RegionUseCase) Update(ctx context.Context, region *entity.Region) error {
	return uc.repo.Update(ctx, region)
}

func (uc *RegionUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *RegionUseCase) ListByCountry(ctx context.Context, countryCode string) ([]entity.Region, error) {
	return uc.repo.ListByCountry(ctx, countryCode)
}
