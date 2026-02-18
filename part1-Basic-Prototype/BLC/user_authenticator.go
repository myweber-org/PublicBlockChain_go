package middleware

import (
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

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

        tokenStr := parts[1]
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        r = r.WithContext(context.WithValue(r.Context(), "user_id", claims.UserID))
        r = r.WithContext(context.WithValue(r.Context(), "user_role", claims.Role))
        
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

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		userID, err := validateToken(token)
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

func validateToken(token string) (string, error) {
	// Placeholder for actual token validation logic
	// In real implementation, parse JWT and verify signature
	if token == "" {
		return "", http.ErrNoCookie
	}
	return "user123", nil
}