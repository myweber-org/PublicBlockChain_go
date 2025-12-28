
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	UserAgent string
	IP        string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		activity := ActivityLog{
			Timestamp: start,
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IP:        r.RemoteAddr,
		}

		log.Printf("Activity: %s %s from %s (%s) at %s",
			activity.Method,
			activity.Path,
			activity.IP,
			activity.UserAgent,
			activity.Timestamp.Format(time.RFC3339),
		)

		next.ServeHTTP(w, r)
	})
}