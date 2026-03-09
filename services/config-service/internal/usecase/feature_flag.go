package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/domain/repository"
)

type FeatureFlagUseCase struct {
	repo repository.FeatureFlagRepository
}

func NewFeatureFlagUseCase(repo repository.FeatureFlagRepository) *FeatureFlagUseCase {
	return &FeatureFlagUseCase{repo: repo}
}

func (uc *FeatureFlagUseCase) Create(ctx context.Context, flag *entity.FeatureFlag) error {
	if flag.CountryCode == "" || flag.FlagKey == "" {
		return fmt.Errorf("%w: country_code and flag_key required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, flag)
}

func (uc *FeatureFlagUseCase) GetByID(ctx context.Context, id string) (*entity.FeatureFlag, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *FeatureFlagUseCase) Update(ctx context.Context, flag *entity.FeatureFlag) error {
	return uc.repo.Update(ctx, flag)
}

func (uc *FeatureFlagUseCase) Toggle(ctx context.Context, id string) (*entity.FeatureFlag, error) {
	flag, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	flag.Enabled = !flag.Enabled
	if err := uc.repo.Update(ctx, flag); err != nil {
		return nil, err
	}
	return flag, nil
}

func (uc *FeatureFlagUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *FeatureFlagUseCase) List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.FeatureFlag], error) {
	return uc.repo.List(ctx, countryCode, page)
}

func (uc *FeatureFlagUseCase) IsEnabled(ctx context.Context, countryCode, flagKey string) (bool, error) {
	return uc.repo.IsEnabled(ctx, countryCode, flagKey)
}
