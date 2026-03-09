package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/domain/repository"
)

type SystemConfigUseCase struct {
	repo repository.SystemConfigRepository
}

func NewSystemConfigUseCase(repo repository.SystemConfigRepository) *SystemConfigUseCase {
	return &SystemConfigUseCase{repo: repo}
}

func (uc *SystemConfigUseCase) Create(ctx context.Context, cfg *entity.SystemConfig) error {
	if cfg.CountryCode == "" || cfg.ConfigKey == "" {
		return fmt.Errorf("%w: country_code and config_key required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, cfg)
}

func (uc *SystemConfigUseCase) GetByID(ctx context.Context, id string) (*entity.SystemConfig, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *SystemConfigUseCase) GetByKey(ctx context.Context, countryCode, key string) (*entity.SystemConfig, error) {
	return uc.repo.GetByKey(ctx, countryCode, key)
}

func (uc *SystemConfigUseCase) Update(ctx context.Context, cfg *entity.SystemConfig) error {
	if cfg.ID == "" {
		return fmt.Errorf("%w: id required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Update(ctx, cfg)
}

func (uc *SystemConfigUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *SystemConfigUseCase) List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.SystemConfig], error) {
	return uc.repo.List(ctx, countryCode, page)
}
