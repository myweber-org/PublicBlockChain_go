package middleware

import (
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	secretKey []byte
}

func NewUserAuthenticator(secretKey string) *UserAuthenticator {
	return &UserAuthenticator{
		secretKey: []byte(secretKey),
	}
}

func (ua *UserAuthenticator) Authenticate(next http.Handler) http.Handler {
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
		claims, err := validateJWTToken(tokenString, ua.secretKey)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func validateJWTToken(tokenString string, secretKey []byte) (*TokenClaims, error) {
	// JWT validation implementation would go here
	// This is a simplified placeholder
	return &TokenClaims{
		UserID: "sample-user-id",
		Role:   "user",
	}, nil
}