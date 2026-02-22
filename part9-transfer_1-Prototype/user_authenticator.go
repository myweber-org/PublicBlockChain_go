package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

func Authenticate(next http.Handler) http.Handler {
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
		userID, err := validateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(token string) (string, error) {
	// Implementation would verify JWT signature and extract claims
	// This is a simplified placeholder
	if token == "" {
		return "", http.ErrNoCookie
	}
	return "user123", nil
}

func GetUserID(ctx context.Context) string {
	if val, ok := ctx.Value(userIDKey).(string); ok {
		return val
	}
	return ""
}package middleware

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
	
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid token format")
	}
	
	return true, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
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
		
		valid, err := a.ValidateToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Token validation error: %v", err), http.StatusUnauthorized)
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
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

func Authenticate(next http.Handler) http.Handler {
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

		tokenString := parts[1]
		userID, err := validateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func validateToken(tokenString string) (string, error) {
	// This is a placeholder for actual JWT validation logic
	// In production, use a proper JWT library like github.com/golang-jwt/jwt
	// For this example, we'll simulate validation by checking a prefix
	if !strings.HasPrefix(tokenString, "valid_") {
		return "", http.ErrAbortHandler
	}
	return strings.TrimPrefix(tokenString, "valid_"), nil
}