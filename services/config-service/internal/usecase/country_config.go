package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/domain/repository"
)

type CountryConfigUseCase struct {
	repo repository.CountryConfigRepository
}

func NewCountryConfigUseCase(repo repository.CountryConfigRepository) *CountryConfigUseCase {
	return &CountryConfigUseCase{repo: repo}
}

func (uc *CountryConfigUseCase) Create(ctx context.Context, cfg *entity.CountryConfig) error {
	if cfg.CountryCode == "" || cfg.CountryName == "" {
		return fmt.Errorf("%w: country_code and country_name required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, cfg)
}

func (uc *CountryConfigUseCase) GetByCode(ctx context.Context, code string) (*entity.CountryConfig, error) {
	return uc.repo.GetByCode(ctx, code)
}

func (uc *CountryConfigUseCase) Update(ctx context.Context, cfg *entity.CountryConfig) error {
	return uc.repo.Update(ctx, cfg)
}

func (uc *CountryConfigUseCase) List(ctx context.Context) ([]entity.CountryConfig, error) {
	return uc.repo.List(ctx)
}

func (uc *CountryConfigUseCase) ListActive(ctx context.Context) ([]entity.CountryConfig, error) {
	return uc.repo.ListActive(ctx)
}
