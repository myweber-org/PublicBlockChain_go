
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	IPAddress string
	Method    string
	Path      string
	Timestamp time.Time
	UserAgent string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserIDFromToken(authHeader)
		}

		activity := ActivityLog{
			UserID:    userID,
			IPAddress: r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
			Timestamp: start,
			UserAgent: r.UserAgent(),
		}

		logActivity(activity)

		next.ServeHTTP(w, r)
	})
}

func extractUserIDFromToken(token string) string {
	return "user_" + token[:8]
}

func logActivity(activity ActivityLog) {
	log.Printf("ACTIVITY: %s %s %s %s %s",
		activity.Timestamp.Format(time.RFC3339),
		activity.UserID,
		activity.IPAddress,
		activity.Method,
		activity.Path)
}