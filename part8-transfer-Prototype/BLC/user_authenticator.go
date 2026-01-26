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
	
	if !strings.HasPrefix(token, "Bearer ") {
		return false
	}
	
	claims, err := parseJWT(strings.TrimPrefix(token, "Bearer "), a.secretKey)
	if err != nil {
		return false
	}
	
	return claims.Valid()
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		
		if !a.ValidateToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}