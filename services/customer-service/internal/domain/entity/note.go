package entity

import "time"

type CustomerNote struct {
	ID         string
	CustomerID string
	Note       string
	CreatedBy  string
	CreatedAt  time.Time
}
