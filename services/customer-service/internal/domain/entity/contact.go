package entity

import "time"

type CustomerContact struct {
	ID         string
	CustomerID string
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Mobile     string
	Position   string
	IsPrimary  bool
	IsBilling  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
