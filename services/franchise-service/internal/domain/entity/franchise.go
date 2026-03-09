package entity

import "time"

type Franchise struct {
	ID                string
	CountryCode       string
	Name              string
	Code              string
	ContactName       string
	Email             string
	Phone             string
	AddressLine1      string
	AddressLine2      string
	City              string
	State             string
	Postcode          string
	IsActive          bool
	CommissionRate    float64
	ParentFranchiseID *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type FranchiseSetting struct {
	ID           string
	FranchiseID  string
	SettingKey   string
	SettingValue string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FranchiseHistory struct {
	ID          string
	FranchiseID string
	Action      string
	ChangedBy   string
	Changes     map[string]any
	CreatedAt   time.Time
}
