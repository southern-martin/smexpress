package dto

import (
	"time"

	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type CreateCountryConfigRequest struct {
	CountryCode    string `json:"country_code"`
	CountryName    string `json:"country_name"`
	CurrencyCode   string `json:"currency_code"`
	CurrencySymbol string `json:"currency_symbol"`
	Timezone       string `json:"timezone"`
	DateFormat     string `json:"date_format"`
	WeightUnit     string `json:"weight_unit"`
	DimensionUnit  string `json:"dimension_unit"`
	Locale         string `json:"locale"`
}

type UpdateCountryConfigRequest struct {
	CountryName    string `json:"country_name"`
	CurrencyCode   string `json:"currency_code"`
	CurrencySymbol string `json:"currency_symbol"`
	Timezone       string `json:"timezone"`
	DateFormat     string `json:"date_format"`
	WeightUnit     string `json:"weight_unit"`
	DimensionUnit  string `json:"dimension_unit"`
	Locale         string `json:"locale"`
	IsActive       *bool  `json:"is_active,omitempty"`
}

type CountryConfigResponse struct {
	ID             string    `json:"id"`
	CountryCode    string    `json:"country_code"`
	CountryName    string    `json:"country_name"`
	CurrencyCode   string    `json:"currency_code"`
	CurrencySymbol string    `json:"currency_symbol"`
	Timezone       string    `json:"timezone"`
	DateFormat     string    `json:"date_format"`
	WeightUnit     string    `json:"weight_unit"`
	DimensionUnit  string    `json:"dimension_unit"`
	Locale         string    `json:"locale"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func CountryConfigFromEntity(e *entity.CountryConfig) CountryConfigResponse {
	return CountryConfigResponse{
		ID:             e.ID,
		CountryCode:    e.CountryCode,
		CountryName:    e.CountryName,
		CurrencyCode:   e.CurrencyCode,
		CurrencySymbol: e.CurrencySymbol,
		Timezone:       e.Timezone,
		DateFormat:     e.DateFormat,
		WeightUnit:     e.WeightUnit,
		DimensionUnit:  e.DimensionUnit,
		Locale:         e.Locale,
		IsActive:       e.IsActive,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
