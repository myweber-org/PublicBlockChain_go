
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
	IPAddress string
	Duration  time.Duration
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		activity := ActivityLog{
			Timestamp: time.Now().UTC(),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
			Duration:  duration,
		}

		log.Printf("ACTIVITY: %s %s | IP: %s | Agent: %s | Duration: %v | Status: %d",
			activity.Method,
			activity.Path,
			activity.IPAddress,
			activity.UserAgent,
			activity.Duration,
			lw.statusCode,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}