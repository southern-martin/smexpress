package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/domain/repository"
)

type HolidayUseCase struct {
	repo repository.HolidayRepository
}

func NewHolidayUseCase(repo repository.HolidayRepository) *HolidayUseCase {
	return &HolidayUseCase{repo: repo}
}

func (uc *HolidayUseCase) Create(ctx context.Context, h *entity.Holiday) error {
	if h.CountryCode == "" || h.Name == "" {
		return fmt.Errorf("%w: country_code and name required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, h)
}

func (uc *HolidayUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *HolidayUseCase) List(ctx context.Context, countryCode string, year int, page db.Page) (db.PagedResult[entity.Holiday], error) {
	return uc.repo.List(ctx, countryCode, year, page)
}
