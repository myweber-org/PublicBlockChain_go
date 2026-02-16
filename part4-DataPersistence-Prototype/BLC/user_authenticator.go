package middleware

import (
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey string
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{secretKey: secretKey}
}

func (a *Authenticator) ValidateToken(tokenString string) (bool, error) {
	if strings.TrimSpace(tokenString) == "" {
		return false, nil
	}
	
	// In production, implement proper JWT validation
	// using libraries like github.com/golang-jwt/jwt
	// This is a simplified example
	return tokenString == a.secretKey, nil
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
		if err != nil {
			http.Error(w, "Token validation error", http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}