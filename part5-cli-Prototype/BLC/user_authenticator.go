package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey string
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{secretKey: secretKey}
}

func (a *Authenticator) ValidateToken(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("empty token")
	}
	
	// Simulated token validation
	if strings.HasPrefix(token, "valid_") && len(token) > 10 {
		return true, nil
	}
	
	return false, fmt.Errorf("invalid token format")
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		valid, err := a.ValidateToken(token)
		if err != nil || !valid {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}