package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/address-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/address-service/internal/domain/errors"
	"github.com/smexpress/services/address-service/internal/domain/repository"
)

type PostcodeUseCase struct {
	repo repository.PostcodeRepository
}

func NewPostcodeUseCase(repo repository.PostcodeRepository) *PostcodeUseCase {
	return &PostcodeUseCase{repo: repo}
}

func (uc *PostcodeUseCase) Search(ctx context.Context, countryCode, query string, page db.Page) (db.PagedResult[entity.Postcode], error) {
	if countryCode == "" {
		return db.PagedResult[entity.Postcode]{}, fmt.Errorf("%w: country_code required", domainerr.ErrInvalidInput)
	}
	if query == "" {
		return db.PagedResult[entity.Postcode]{}, fmt.Errorf("%w: search query required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Search(ctx, countryCode, query, page)
}

func (uc *PostcodeUseCase) Lookup(ctx context.Context, countryCode, postcode string) ([]entity.Postcode, error) {
	if countryCode == "" || postcode == "" {
		return nil, fmt.Errorf("%w: country_code and postcode required", domainerr.ErrInvalidInput)
	}
	return uc.repo.GetByPostcode(ctx, countryCode, postcode)
}

func (uc *PostcodeUseCase) Create(ctx context.Context, p *entity.Postcode) error {
	if p.CountryCode == "" || p.Postcode == "" {
		return fmt.Errorf("%w: country_code and postcode required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, p)
}

func (uc *PostcodeUseCase) BulkCreate(ctx context.Context, postcodes []entity.Postcode) (int64, error) {
	if len(postcodes) == 0 {
		return 0, fmt.Errorf("%w: no postcodes provided", domainerr.ErrInvalidInput)
	}
	return uc.repo.BulkCreate(ctx, postcodes)
}
