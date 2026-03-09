package errors

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account locked")
	ErrAccountInactive    = errors.New("account inactive")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrTokenExpired       = errors.New("token expired")
)
