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

type WithdrawalUseCase struct {
	withdrawalRepo repository.WithdrawalRepository
	ledgerUC       *LedgerUseCase
}

func NewWithdrawalUseCase(withdrawalRepo repository.WithdrawalRepository, ledgerUC *LedgerUseCase) *WithdrawalUseCase {
	return &WithdrawalUseCase{withdrawalRepo: withdrawalRepo, ledgerUC: ledgerUC}
}

func (uc *WithdrawalUseCase) Request(ctx context.Context, w *entity.FranchiseWithdrawal) error {
	if w.Amount <= 0 {
		return fmt.Errorf("%w: amount must be positive", domainerr.ErrInvalidInput)
	}

	ledger, err := uc.ledgerUC.GetBalance(ctx, w.FranchiseID)
	if err != nil {
		return err
	}
	if ledger.Balance < w.Amount {
		return fmt.Errorf("%w: balance %.2f, requested %.2f", domainerr.ErrInsufficientBalance, ledger.Balance, w.Amount)
	}

	w.Status = "pending"
	w.RequestedAt = time.Now()
	return uc.withdrawalRepo.Create(ctx, w)
}

func (uc *WithdrawalUseCase) Approve(ctx context.Context, id, approvedBy string) error {
	w, err := uc.withdrawalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := w.Approve(approvedBy); err != nil {
		return err
	}
	return uc.withdrawalRepo.Update(ctx, w)
}

func (uc *WithdrawalUseCase) Reject(ctx context.Context, id, approvedBy string) error {
	w, err := uc.withdrawalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := w.Reject(approvedBy); err != nil {
		return err
	}
	return uc.withdrawalRepo.Update(ctx, w)
}

func (uc *WithdrawalUseCase) Process(ctx context.Context, id string) error {
	w, err := uc.withdrawalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := w.Process(); err != nil {
		return err
	}

	if err := uc.ledgerUC.Debit(ctx, w.FranchiseID, w.Amount, "Withdrawal processed", "withdrawal", w.ID); err != nil {
		return err
	}
	return uc.withdrawalRepo.Update(ctx, w)
}

func (uc *WithdrawalUseCase) GetByID(ctx context.Context, id string) (*entity.FranchiseWithdrawal, error) {
	return uc.withdrawalRepo.GetByID(ctx, id)
}

func (uc *WithdrawalUseCase) ListByFranchise(ctx context.Context, franchiseID string, page db.Page) (db.PagedResult[entity.FranchiseWithdrawal], error) {
	return uc.withdrawalRepo.ListByFranchise(ctx, franchiseID, page)
}

func (uc *WithdrawalUseCase) ListPending(ctx context.Context, page db.Page) (db.PagedResult[entity.FranchiseWithdrawal], error) {
	return uc.withdrawalRepo.ListPending(ctx, page)
}
