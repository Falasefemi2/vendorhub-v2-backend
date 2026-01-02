package utils

import "context"

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", ErrUnauthorized
	}
	return userID, nil
}

func GetRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value(RoleKey).(string)
	if !ok || role == "" {
		return "", ErrUnauthorized
	}
	return role, nil
}
