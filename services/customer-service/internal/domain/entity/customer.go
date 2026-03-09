package entity

import "time"

type Customer struct {
	ID             string
	CountryCode    string
	FranchiseID    string
	CompanyName    string
	TradingName    string
	AccountNumber  string
	ABN            string
	Email          string
	Phone          string
	Website        string
	CreditLimit    float64
	CreditBalance  float64
	PaymentTerms   int
	IsActive       bool
	IsCreditHold   bool
	Notes          string
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
