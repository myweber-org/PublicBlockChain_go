package middleware

import (
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
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
                return []byte(secretKey), nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
                return
            }

            r.Header.Set("X-Username", claims.Username)
            r.Header.Set("X-User-Role", claims.Role)
            next.ServeHTTP(w, r)
        })
    }
}