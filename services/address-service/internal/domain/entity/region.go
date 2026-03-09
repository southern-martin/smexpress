package entity

import "time"

type Region struct {
	ID             string
	CountryCode    string
	Name           string
	Code           string
	ParentRegionID *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
