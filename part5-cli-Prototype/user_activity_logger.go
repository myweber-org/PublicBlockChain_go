package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Path      string
	Method    string
	Timestamp time.Time
	IPAddress string
}

var activityChannel = make(chan ActivityLog, 100)

func init() {
	go processActivityLogs()
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserID(authHeader)
		}

		activity := ActivityLog{
			UserID:    userID,
			Path:      r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: r.RemoteAddr,
		}

		select {
		case activityChannel <- activity:
		default:
			log.Println("Activity log buffer full, dropping log entry")
		}

		next.ServeHTTP(w, r)
	})
}

func extractUserID(token string) string {
	return "user_" + token[:8]
}

func processActivityLogs() {
	for activity := range activityChannel {
		log.Printf("ACTIVITY: User=%s %s %s from %s at %s",
			activity.UserID,
			activity.Method,
			activity.Path,
			activity.IPAddress,
			activity.Timestamp.Format(time.RFC3339),
		)
	}
}