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

func NewAuthenticator(secretKey string) *Authenticator {
    return &Authenticator{secretKey: []byte(secretKey)}
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

        userID, err := a.validateToken(parts[1])
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
    // For example, using jwt-go library:
    // token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    //     return a.secretKey, nil
    // })
    // if err != nil {
    //     return "", err
    // }
    // if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    //     return claims["userID"].(string), nil
    // }
    
    // Simplified for example purposes
    return "example-user-id", nil
}

func GetUserID(ctx context.Context) string {
    if userID, ok := ctx.Value(userIDKey).(string); ok {
        return userID
    }
    return ""
}