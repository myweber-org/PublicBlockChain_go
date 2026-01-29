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

func (a *Authenticator) ValidateToken(token string) bool {
    if token == "" {
        return false
    }
    
    // Simulated token validation logic
    // In real implementation, use proper JWT validation library
    return strings.HasPrefix(token, "valid_")
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if !a.ValidateToken(token) {
            http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}