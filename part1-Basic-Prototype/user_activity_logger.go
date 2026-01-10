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
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := "anonymous"
		if auth := r.Header.Get("Authorization"); auth != "" {
			userID = extractUserID(auth)
		}

		activity := ActivityLog{
			UserID:    userID,
			IPAddress: r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
			Timestamp: start,
		}

		logActivity(activity)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:min(len(authHeader), 20)] + "..."
	}
	return "authenticated"
}

func logActivity(activity ActivityLog) {
	log.Printf("ACTIVITY: User=%s IP=%s %s %s at %s",
		activity.UserID,
		activity.IPAddress,
		activity.Method,
		activity.Path,
		activity.Timestamp.Format(time.RFC3339))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}