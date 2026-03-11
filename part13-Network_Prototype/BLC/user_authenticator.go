package middleware

import (
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey []byte
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{
		secretKey: []byte(secretKey),
	}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		if !a.validateToken(tokenString) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authenticator) validateToken(tokenString string) bool {
	// Simplified token validation logic
	// In production, use proper JWT library like github.com/golang-jwt/jwt
	return len(tokenString) > 10 && strings.HasPrefix(tokenString, "valid_")
}