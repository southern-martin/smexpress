package entity

import "time"

type Postcode struct {
	ID          string
	CountryCode string
	Postcode    string
	Suburb      string
	City        string
	State       string
	StateCode   string
	Latitude    *float64
	Longitude   *float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
