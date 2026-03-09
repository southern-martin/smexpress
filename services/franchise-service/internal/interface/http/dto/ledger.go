package dto

import (
	"time"

	"github.com/smexpress/services/franchise-service/internal/domain/entity"
)

type CreditRequest struct {
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	ReferenceType string  `json:"reference_type"`
	ReferenceID   string  `json:"reference_id"`
}

type LedgerResponse struct {
	ID          string    `json:"id"`
	FranchiseID string    `json:"franchise_id"`
	CountryCode string    `json:"country_code"`
	Currency    string    `json:"currency"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func LedgerFromEntity(l *entity.FranchiseLedger) LedgerResponse {
	return LedgerResponse{
		ID:          l.ID,
		FranchiseID: l.FranchiseID,
		CountryCode: l.CountryCode,
		Currency:    l.Currency,
		Balance:     l.Balance,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

type LedgerEntryResponse struct {
	ID            string    `json:"id"`
	LedgerID      string    `json:"ledger_id"`
	EntryType     string    `json:"entry_type"`
	Amount        float64   `json:"amount"`
	BalanceAfter  float64   `json:"balance_after"`
	Description   string    `json:"description"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func LedgerEntryFromEntity(e *entity.FranchiseLedgerEntry) LedgerEntryResponse {
	return LedgerEntryResponse{
		ID:            e.ID,
		LedgerID:      e.LedgerID,
		EntryType:     e.EntryType,
		Amount:        e.Amount,
		BalanceAfter:  e.BalanceAfter,
		Description:   e.Description,
		ReferenceType: e.ReferenceType,
		ReferenceID:   e.ReferenceID,
		CreatedAt:     e.CreatedAt,
	}
}
