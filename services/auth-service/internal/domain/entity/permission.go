package entity

import "time"

type Permission struct {
	ID          string
	Code        string
	Name        string
	Module      string
	Description string
	CreatedAt   time.Time
}
