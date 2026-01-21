package auth

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

const userIDKey contextKey = "userID"

type Authenticator struct {
    secretKey []byte
}

func NewAuthenticator(secret string) *Authenticator {
    return &Authenticator{secretKey: []byte(secret)}
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

        tokenString := parts[1]
        userID, err := a.validateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (a *Authenticator) validateToken(tokenString string) (string, error) {
    // Token validation logic would go here
    // For this example, we'll use a simple placeholder
    if tokenString == "valid_token_example" {
        return "user123", nil
    }
    return "", http.ErrNoCookie
}

func GetUserID(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDKey).(string)
    return userID, ok
}