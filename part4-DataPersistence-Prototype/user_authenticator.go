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
	if tokenString == "" {
		return false, nil
	}
	
	// Simplified token validation logic
	// In production, use proper JWT library like github.com/golang-jwt/jwt
	expectedPrefix := "Bearer "
	if !strings.HasPrefix(tokenString, expectedPrefix) {
		return false, nil
	}
	
	token := strings.TrimPrefix(tokenString, expectedPrefix)
	return len(token) > 10, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		
		valid, err := a.ValidateToken(token)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}