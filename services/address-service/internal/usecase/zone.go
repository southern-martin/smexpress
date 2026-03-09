package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/address-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/address-service/internal/domain/errors"
	"github.com/smexpress/services/address-service/internal/domain/repository"
)

type ZoneUseCase struct {
	repo repository.ZoneRepository
}

func NewZoneUseCase(repo repository.ZoneRepository) *ZoneUseCase {
	return &ZoneUseCase{repo: repo}
}

func (uc *ZoneUseCase) Create(ctx context.Context, zone *entity.Zone) error {
	if zone.CountryCode == "" || zone.ZoneName == "" || zone.ZoneCode == "" {
		return fmt.Errorf("%w: country_code, zone_name and zone_code required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, zone)
}

func (uc *ZoneUseCase) GetByID(ctx context.Context, id string) (*entity.Zone, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ZoneUseCase) Update(ctx context.Context, zone *entity.Zone) error {
	return uc.repo.Update(ctx, zone)
}

func (uc *ZoneUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *ZoneUseCase) ListByCountry(ctx context.Context, countryCode string) ([]entity.Zone, error) {
	return uc.repo.ListByCountry(ctx, countryCode)
}

func (uc *ZoneUseCase) FindZone(ctx context.Context, countryCode, postcode string) (*entity.Zone, error) {
	if countryCode == "" || postcode == "" {
		return nil, fmt.Errorf("%w: country_code and postcode required", domainerr.ErrInvalidInput)
	}
	return uc.repo.FindZoneForPostcode(ctx, countryCode, postcode)
}

func (uc *ZoneUseCase) SetPostcodes(ctx context.Context, zoneID string, postcodes []entity.ZonePostcode) error {
	return uc.repo.SetPostcodes(ctx, zoneID, postcodes)
}
