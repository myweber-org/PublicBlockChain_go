package middleware

import (
	"net/http"
	"strings"
)

type User struct {
	ID    int
	Roles []string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(token)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if !hasRequiredRole(user, r) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func validateToken(token string) (*User, error) {
	// Token validation logic
	if token == "valid_token_example" {
		return &User{ID: 1, Roles: []string{"admin", "user"}}, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func hasRequiredRole(user *User, r *http.Request) bool {
	requiredRole := r.Header.Get("X-Required-Role")
	if requiredRole == "" {
		return true
	}
	for _, role := range user.Roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}