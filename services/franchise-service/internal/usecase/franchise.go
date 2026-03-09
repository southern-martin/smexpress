package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
	"github.com/smexpress/services/franchise-service/internal/domain/repository"
)

type FranchiseUseCase struct {
	franchiseRepo repository.FranchiseRepository
	ledgerRepo    repository.LedgerRepository
}

func NewFranchiseUseCase(franchiseRepo repository.FranchiseRepository, ledgerRepo repository.LedgerRepository) *FranchiseUseCase {
	return &FranchiseUseCase{franchiseRepo: franchiseRepo, ledgerRepo: ledgerRepo}
}

func (uc *FranchiseUseCase) Create(ctx context.Context, franchise *entity.Franchise) error {
	if franchise.Name == "" || franchise.Code == "" || franchise.CountryCode == "" {
		return fmt.Errorf("%w: name, code and country_code required", domainerr.ErrInvalidInput)
	}

	if err := uc.franchiseRepo.Create(ctx, franchise); err != nil {
		return err
	}

	ledger := &entity.FranchiseLedger{
		FranchiseID: franchise.ID,
		CountryCode: franchise.CountryCode,
		Currency:    "AUD",
		Balance:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return uc.ledgerRepo.Create(ctx, ledger)
}

func (uc *FranchiseUseCase) GetByID(ctx context.Context, id string) (*entity.Franchise, error) {
	return uc.franchiseRepo.GetByID(ctx, id)
}

func (uc *FranchiseUseCase) Update(ctx context.Context, franchise *entity.Franchise) error {
	return uc.franchiseRepo.Update(ctx, franchise)
}

func (uc *FranchiseUseCase) List(ctx context.Context, countryCode string, page db.Page) (db.PagedResult[entity.Franchise], error) {
	return uc.franchiseRepo.List(ctx, countryCode, page)
}

func (uc *FranchiseUseCase) Activate(ctx context.Context, id string) error {
	f, err := uc.franchiseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	f.IsActive = true
	return uc.franchiseRepo.Update(ctx, f)
}

func (uc *FranchiseUseCase) Deactivate(ctx context.Context, id string) error {
	f, err := uc.franchiseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	f.IsActive = false
	return uc.franchiseRepo.Update(ctx, f)
}

func (uc *FranchiseUseCase) GetSettings(ctx context.Context, franchiseID string) ([]entity.FranchiseSetting, error) {
	return uc.franchiseRepo.GetSettings(ctx, franchiseID)
}

func (uc *FranchiseUseCase) SetSetting(ctx context.Context, franchiseID, key, value string) error {
	if key == "" {
		return fmt.Errorf("%w: setting key required", domainerr.ErrInvalidInput)
	}
	return uc.franchiseRepo.SetSetting(ctx, franchiseID, key, value)
}
