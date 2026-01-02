package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func HandleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		WriteError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrInvalidCredentials):
		WriteError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrAccountNotActive):
		WriteError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrInvalidOperation):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrUserNotFound):
		WriteError(w, http.StatusNotFound, err.Error())
	default:
		WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

// ParseFloat64 parses a string to float64
func ParseFloat64(value string) (float64, error) {
	if value == "" {
		return 0, errors.New("value cannot be empty")
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float value: %w", err)
	}

	if f < 0 {
		return 0, errors.New("value cannot be negative")
	}

	return f, nil
}
