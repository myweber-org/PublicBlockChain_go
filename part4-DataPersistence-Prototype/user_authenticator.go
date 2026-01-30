package middleware

import (
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	secretKey string
}

func NewUserAuthenticator(secret string) *UserAuthenticator {
	return &UserAuthenticator{secretKey: secret}
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
		
		valid, userID := ua.ValidateToken(parts[1])
		if !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		
		r.Header.Set("X-User-ID", userID)
		next.ServeHTTP(w, r)
	})
}

func parseJWTToken(token, secret string) (*TokenClaims, error) {
	// Implementation would use JWT library like github.com/golang-jwt/jwt
	// This is a simplified placeholder
	return &TokenClaims{UserID: "sample-user-id"}, nil
}

type TokenClaims struct {
	UserID string
}package auth

import (
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

var jwtKey = []byte("your_secret_key_here")

func GenerateToken(username, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, jwt.ErrSignatureInvalid
    }

    return claims, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
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

        claims, err := ValidateToken(parts[1])
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        r.Header.Set("X-Username", claims.Username)
        r.Header.Set("X-Role", claims.Role)
        next.ServeHTTP(w, r)
    })
}