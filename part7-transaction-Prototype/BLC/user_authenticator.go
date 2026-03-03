package middleware

import (
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey []byte
}

func NewAuthenticator(secret string) *Authenticator {
	return &Authenticator{secretKey: []byte(secret)}
}

func (a *Authenticator) ValidateToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, nil
	}
	
	// Simulate token validation logic
	// In real implementation, this would parse and verify JWT
	if strings.HasPrefix(tokenString, "valid_") {
		return true, nil
	}
	
	return false, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}
		
		valid, err := a.ValidateToken(parts[1])
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
}package middleware

import (
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	secretKey string
}

func NewUserAuthenticator(secretKey string) *UserAuthenticator {
	return &UserAuthenticator{secretKey: secretKey}
}

func (ua *UserAuthenticator) ValidateToken(token string) (bool, string) {
	if token == "" {
		return false, ""
	}
	
	claims, err := parseJWTToken(token, ua.secretKey)
	if err != nil {
		return false, ""
	}
	
	return true, claims.UserID
}

func (ua *UserAuthenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}
		
		token := parts[1]
		valid, userID := ua.ValidateToken(token)
		if !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		
		r.Header.Set("X-User-ID", userID)
		next.ServeHTTP(w, r)
	})
}

func parseJWTToken(token, secretKey string) (*TokenClaims, error) {
	// JWT parsing implementation would go here
	// This is a simplified placeholder
	return &TokenClaims{UserID: "sample-user-id"}, nil
}

type TokenClaims struct {
	UserID string
}