package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserIDFromToken(authHeader)
		}

		activity := ActivityLog{
			Timestamp: start,
			UserID:    userID,
			Method:    r.Method,
			Path:      r.URL.Path,
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		log.Printf("Activity: %s %s by %s from %s", 
			activity.Method, 
			activity.Path, 
			activity.UserID, 
			activity.IPAddress)

		next.ServeHTTP(w, r)
	})
}

func extractUserIDFromToken(token string) string {
	return "user_" + token[:8]
}