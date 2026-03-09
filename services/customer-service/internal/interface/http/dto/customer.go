package dto

import (
	"time"

	"github.com/smexpress/services/customer-service/internal/domain/entity"
)

type CreateCustomerRequest struct {
	CountryCode   string  `json:"country_code"`
	FranchiseID   string  `json:"franchise_id"`
	CompanyName   string  `json:"company_name"`
	TradingName   string  `json:"trading_name"`
	AccountNumber string  `json:"account_number"`
	ABN           string  `json:"abn"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	Website       string  `json:"website"`
	CreditLimit   float64 `json:"credit_limit"`
	PaymentTerms  int     `json:"payment_terms"`
}

type UpdateCustomerRequest struct {
	CompanyName  string   `json:"company_name"`
	TradingName  string   `json:"trading_name"`
	ABN          string   `json:"abn"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	Website      string   `json:"website"`
	CreditLimit  *float64 `json:"credit_limit,omitempty"`
	PaymentTerms *int     `json:"payment_terms,omitempty"`
}

type CustomerResponse struct {
	ID             string    `json:"id"`
	CountryCode    string    `json:"country_code"`
	FranchiseID    string    `json:"franchise_id,omitempty"`
	CompanyName    string    `json:"company_name"`
	TradingName    string    `json:"trading_name"`
	AccountNumber  string    `json:"account_number"`
	ABN            string    `json:"abn"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Website        string    `json:"website"`
	CreditLimit    float64   `json:"credit_limit"`
	CreditBalance  float64   `json:"credit_balance"`
	PaymentTerms   int       `json:"payment_terms"`
	IsActive       bool      `json:"is_active"`
	IsCreditHold   bool      `json:"is_credit_hold"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func CustomerFromEntity(c *entity.Customer) CustomerResponse {
	return CustomerResponse{
		ID:            c.ID,
		CountryCode:   c.CountryCode,
		FranchiseID:   c.FranchiseID,
		CompanyName:   c.CompanyName,
		TradingName:   c.TradingName,
		AccountNumber: c.AccountNumber,
		ABN:           c.ABN,
		Email:         c.Email,
		Phone:         c.Phone,
		Website:       c.Website,
		CreditLimit:   c.CreditLimit,
		CreditBalance: c.CreditBalance,
		PaymentTerms:  c.PaymentTerms,
		IsActive:      c.IsActive,
		IsCreditHold:  c.IsCreditHold,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}
