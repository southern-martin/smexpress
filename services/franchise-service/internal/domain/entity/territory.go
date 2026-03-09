package entity

import "time"

type Territory struct {
	ID          string
	FranchiseID string
	CountryCode string
	Name        string
	PostcodeFrom string
	PostcodeTo   string
	State       string
	IsExclusive bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
