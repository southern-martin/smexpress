package entity

import "time"

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	RevokedAt *time.Time
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}
