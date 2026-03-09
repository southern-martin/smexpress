package entity

import "time"

type Role struct {
	ID          string
	CountryCode string
	Name        string
	DisplayName string
	Description string
	IsSystem    bool
	Permissions []Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
