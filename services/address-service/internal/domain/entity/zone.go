package entity

import "time"

type Zone struct {
	ID          string
	CountryCode string
	ZoneName    string
	ZoneCode    string
	Description string
	Postcodes   []ZonePostcode
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ZonePostcode struct {
	ID           string
	ZoneID       string
	PostcodeFrom string
	PostcodeTo   string
	CreatedAt    time.Time
}
