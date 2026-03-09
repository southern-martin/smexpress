package entity

import "time"

type SystemConfig struct {
	ID          string
	CountryCode string
	ConfigKey   string
	ConfigValue string
	Description string
	DataType    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
