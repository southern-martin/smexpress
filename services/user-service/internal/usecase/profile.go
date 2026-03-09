package usecase

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/user-service/internal/domain/entity"
	"github.com/smexpress/services/user-service/internal/domain/repository"
)

type ProfileUseCase struct {
	profileRepo repository.UserProfileRepository
	prefRepo    repository.UserPreferenceRepository
}

func NewProfileUseCase(profileRepo repository.UserProfileRepository, prefRepo repository.UserPreferenceRepository) *ProfileUseCase {
	return &ProfileUseCase{profileRepo: profileRepo, prefRepo: prefRepo}
}

func (uc *ProfileUseCase) Create(ctx context.Context, profile *entity.UserProfile) error {
	return uc.profileRepo.Create(ctx, profile)
}

func (uc *ProfileUseCase) GetByUserID(ctx context.Context, userID string) (*entity.UserProfile, error) {
	return uc.profileRepo.GetByUserID(ctx, userID)
}

func (uc *ProfileUseCase) Update(ctx context.Context, profile *entity.UserProfile) error {
	return uc.profileRepo.Update(ctx, profile)
}

func (uc *ProfileUseCase) List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.UserProfile], error) {
	return uc.profileRepo.List(ctx, countryCode, page)
}

func (uc *ProfileUseCase) GetPreferences(ctx context.Context, userID string) ([]entity.UserPreference, error) {
	return uc.prefRepo.ListByUser(ctx, userID)
}

func (uc *ProfileUseCase) SetPreference(ctx context.Context, userID, key, value string) error {
	return uc.prefRepo.Set(ctx, userID, key, value)
}

func (uc *ProfileUseCase) DeletePreference(ctx context.Context, userID, key string) error {
	return uc.prefRepo.Delete(ctx, userID, key)
}
