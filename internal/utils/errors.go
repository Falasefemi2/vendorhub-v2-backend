package utils

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountNotActive   = errors.New("account not active")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUnauthorized       = errors.New("unauthorized user")
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrWeakPassword       = errors.New("password too weak")
	ErrUserNotFound       = errors.New("user not found")
)
