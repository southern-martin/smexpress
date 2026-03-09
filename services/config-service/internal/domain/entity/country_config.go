package entity

import "time"

type CountryConfig struct {
	ID             string
	CountryCode    string
	CountryName    string
	CurrencyCode   string
	CurrencySymbol string
	Timezone       string
	DateFormat     string
	WeightUnit     string
	DimensionUnit  string
	Locale         string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
