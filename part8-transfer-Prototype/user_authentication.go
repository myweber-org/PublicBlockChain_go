package middleware

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	secretKey []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secretKey: []byte(secret)}
}

func (am *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := am.parseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("User-ID", claims.UserID)
		r.Header.Set("User-Role", claims.Role)
		
		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) parseToken(tokenString string) (*TokenClaims, error) {
	// Token parsing implementation would go here
	// This is a simplified placeholder
	return &TokenClaims{
		UserID: "sample-user-id",
		Role:   "user",
	}, nil
}

type TokenClaims struct {
	UserID string
	Role   string
}