package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const maxFailedAttempts = 5

type User struct {
	ID                  string
	CountryCode         string
	Email               string
	PasswordHash        string
	FirstName           string
	LastName            string
	IsActive            bool
	IsLocked            bool
	FailedLoginAttempts int
	LastLoginAt         *time.Time
	PasswordChangedAt   *time.Time
	FranchiseID         string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Roles               []Role
}

func (u *User) CheckPassword(plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain))
}

func (u *User) IncrementFailedAttempts() {
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= maxFailedAttempts {
		u.IsLocked = true
	}
}

func (u *User) ResetFailedAttempts() {
	u.FailedLoginAttempts = 0
}

func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
