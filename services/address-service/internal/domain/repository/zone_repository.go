package repository

import (
	"context"

	"github.com/smexpress/services/address-service/internal/domain/entity"
)

type ZoneRepository interface {
	Create(ctx context.Context, zone *entity.Zone) error
	GetByID(ctx context.Context, id string) (*entity.Zone, error)
	Update(ctx context.Context, zone *entity.Zone) error
	Delete(ctx context.Context, id string) error
	ListByCountry(ctx context.Context, countryCode string) ([]entity.Zone, error)
	FindZoneForPostcode(ctx context.Context, countryCode, postcode string) (*entity.Zone, error)
	SetPostcodes(ctx context.Context, zoneID string, postcodes []entity.ZonePostcode) error
}
