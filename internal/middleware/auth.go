package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/falasefemi2/vendorhub/internal/utils"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}
		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		ctx := context.WithValue(r.Context(), utils.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, utils.RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleValue := r.Context().Value(utils.RoleKey)
		role, ok := roleValue.(string)
		if !ok || role != "admin" {
			utils.WriteError(w, http.StatusForbidden, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
