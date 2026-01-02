// Package utils
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const passwordMinLength = 8

func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	return err == nil
}

func ValidatePassword(password string) error {
	if len(password) < passwordMinLength {
		return ErrWeakPassword
	}
	return nil
}
