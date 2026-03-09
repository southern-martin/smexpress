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

type LedgerUseCase struct {
	repo repository.LedgerRepository
}

func NewLedgerUseCase(repo repository.LedgerRepository) *LedgerUseCase {
	return &LedgerUseCase{repo: repo}
}

func (uc *LedgerUseCase) GetBalance(ctx context.Context, franchiseID string) (*entity.FranchiseLedger, error) {
	return uc.repo.GetByFranchiseID(ctx, franchiseID)
}

func (uc *LedgerUseCase) ListEntries(ctx context.Context, franchiseID string, page db.Page) (db.PagedResult[entity.FranchiseLedgerEntry], error) {
	ledger, err := uc.repo.GetByFranchiseID(ctx, franchiseID)
	if err != nil {
		return db.PagedResult[entity.FranchiseLedgerEntry]{}, err
	}
	return uc.repo.ListEntries(ctx, ledger.ID, page)
}

func (uc *LedgerUseCase) Credit(ctx context.Context, franchiseID string, amount float64, description, refType, refID string) error {
	if amount <= 0 {
		return fmt.Errorf("%w: amount must be positive", domainerr.ErrInvalidInput)
	}

	ledger, err := uc.repo.GetByFranchiseID(ctx, franchiseID)
	if err != nil {
		return err
	}

	newBalance := ledger.Balance + amount
	entry := &entity.FranchiseLedgerEntry{
		LedgerID:      ledger.ID,
		EntryType:     "credit",
		Amount:        amount,
		BalanceAfter:  newBalance,
		Description:   description,
		ReferenceType: refType,
		ReferenceID:   refID,
		CreatedAt:     time.Now(),
	}
	return uc.repo.AddEntry(ctx, ledger.ID, entry)
}

func (uc *LedgerUseCase) Debit(ctx context.Context, franchiseID string, amount float64, description, refType, refID string) error {
	if amount <= 0 {
		return fmt.Errorf("%w: amount must be positive", domainerr.ErrInvalidInput)
	}

	ledger, err := uc.repo.GetByFranchiseID(ctx, franchiseID)
	if err != nil {
		return err
	}

	if ledger.Balance < amount {
		return fmt.Errorf("%w: balance %.2f, requested %.2f", domainerr.ErrInsufficientBalance, ledger.Balance, amount)
	}

	newBalance := ledger.Balance - amount
	entry := &entity.FranchiseLedgerEntry{
		LedgerID:      ledger.ID,
		EntryType:     "debit",
		Amount:        -amount,
		BalanceAfter:  newBalance,
		Description:   description,
		ReferenceType: refType,
		ReferenceID:   refID,
		CreatedAt:     time.Now(),
	}
	return uc.repo.AddEntry(ctx, ledger.ID, entry)
}
