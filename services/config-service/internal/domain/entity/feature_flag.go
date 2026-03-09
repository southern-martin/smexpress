package entity

import "time"

type FeatureFlag struct {
	ID          string
	CountryCode string
	FlagKey     string
	Enabled     bool
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
