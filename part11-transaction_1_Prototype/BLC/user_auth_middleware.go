package middleware

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	allowedTokens map[string]bool
}

func NewAuthMiddleware(validTokens []string) *AuthMiddleware {
	tokenMap := make(map[string]bool)
	for _, token := range validTokens {
		tokenMap[token] = true
	}
	return &AuthMiddleware{allowedTokens: tokenMap}
}

func (am *AuthMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(parts[1])
		if !am.allowedTokens[token] {
			http.Error(w, "Invalid authentication token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}