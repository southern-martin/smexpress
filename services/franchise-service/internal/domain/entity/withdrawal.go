package entity

import (
	"fmt"
	"time"

	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
)

type FranchiseWithdrawal struct {
	ID                string
	FranchiseID       string
	CountryCode       string
	Amount            float64
	Currency          string
	Status            string
	RequestedBy       string
	ApprovedBy        *string
	BankAccountName   string
	BankAccountNumber string
	BankBSB           string
	Notes             string
	RequestedAt       time.Time
	ProcessedAt       *time.Time
}

func (w *FranchiseWithdrawal) Approve(approvedBy string) error {
	if w.Status != "pending" {
		return fmt.Errorf("%w: withdrawal is not pending", domainerr.ErrInvalidInput)
	}
	w.Status = "approved"
	w.ApprovedBy = &approvedBy
	return nil
}

func (w *FranchiseWithdrawal) Reject(approvedBy string) error {
	if w.Status != "pending" {
		return fmt.Errorf("%w: withdrawal is not pending", domainerr.ErrInvalidInput)
	}
	w.Status = "rejected"
	w.ApprovedBy = &approvedBy
	return nil
}

func (w *FranchiseWithdrawal) Process() error {
	if w.Status != "approved" {
		return fmt.Errorf("%w: withdrawal is not approved", domainerr.ErrInvalidInput)
	}
	w.Status = "processed"
	now := time.Now()
	w.ProcessedAt = &now
	return nil
}
