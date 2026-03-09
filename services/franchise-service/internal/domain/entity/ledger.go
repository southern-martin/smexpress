package entity

import "time"

type FranchiseLedger struct {
	ID          string
	FranchiseID string
	CountryCode string
	Currency    string
	Balance     float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FranchiseLedgerEntry struct {
	ID            string
	LedgerID      string
	EntryType     string
	Amount        float64
	BalanceAfter  float64
	Description   string
	ReferenceType string
	ReferenceID   string
	CreatedAt     time.Time
}
