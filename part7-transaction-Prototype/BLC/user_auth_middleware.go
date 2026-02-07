package middleware

import (
	"net/http"
	"strings"
)

type User struct {
	ID    string
	Email string
	Role  string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(token string) (*User, error) {
	// In production, implement proper JWT validation
	// This is a simplified example
	if token == "valid_token_example" {
		return &User{
			ID:    "user_123",
			Email: "user@example.com",
			Role:  "member",
		}, nil
	}
	return nil, fmt.Errorf("invalid token")
}