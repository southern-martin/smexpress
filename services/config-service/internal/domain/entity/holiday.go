package entity

import "time"

type Holiday struct {
	ID          string
	CountryCode string
	HolidayDate time.Time
	Name        string
	IsRecurring bool
	CreatedAt   time.Time
}
