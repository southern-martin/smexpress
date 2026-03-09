package entity

import "time"

type CustomerAddress struct {
	ID           string
	CustomerID   string
	AddressType  string
	CompanyName  string
	ContactName  string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Postcode     string
	CountryCode  string
	Phone        string
	Email        string
	IsDefault    bool
	Instructions string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
