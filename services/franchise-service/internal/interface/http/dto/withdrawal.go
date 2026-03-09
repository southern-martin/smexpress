package dto

import (
	"time"

	"github.com/smexpress/services/franchise-service/internal/domain/entity"
)

type CreateWithdrawalRequest struct {
	FranchiseID       string  `json:"franchise_id"`
	CountryCode       string  `json:"country_code"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	BankAccountName   string  `json:"bank_account_name"`
	BankAccountNumber string  `json:"bank_account_number"`
	BankBSB           string  `json:"bank_bsb"`
	Notes             string  `json:"notes"`
}

type WithdrawalResponse struct {
	ID                string     `json:"id"`
	FranchiseID       string     `json:"franchise_id"`
	CountryCode       string     `json:"country_code"`
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	Status            string     `json:"status"`
	RequestedBy       string     `json:"requested_by"`
	ApprovedBy        *string    `json:"approved_by,omitempty"`
	BankAccountName   string     `json:"bank_account_name"`
	BankAccountNumber string     `json:"bank_account_number"`
	BankBSB           string     `json:"bank_bsb"`
	Notes             string     `json:"notes"`
	RequestedAt       time.Time  `json:"requested_at"`
	ProcessedAt       *time.Time `json:"processed_at,omitempty"`
}

func WithdrawalFromEntity(w *entity.FranchiseWithdrawal) WithdrawalResponse {
	return WithdrawalResponse{
		ID:                w.ID,
		FranchiseID:       w.FranchiseID,
		CountryCode:       w.CountryCode,
		Amount:            w.Amount,
		Currency:          w.Currency,
		Status:            w.Status,
		RequestedBy:       w.RequestedBy,
		ApprovedBy:        w.ApprovedBy,
		BankAccountName:   w.BankAccountName,
		BankAccountNumber: w.BankAccountNumber,
		BankBSB:           w.BankBSB,
		Notes:             w.Notes,
		RequestedAt:       w.RequestedAt,
		ProcessedAt:       w.ProcessedAt,
	}
}
