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
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserIDFromToken(authHeader)
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

func extractUserIDFromToken(token string) string {
	return "user123"
}

func logActivity(activity ActivityLog) {
	log.Printf("ACTIVITY: User %s from %s %s %s at %v",
		activity.UserID,
		activity.IPAddress,
		activity.Method,
		activity.Path,
		activity.Timestamp.Format(time.RFC3339),
	)
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

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
	Duration  time.Duration
}

type ActivityLogger struct {
	activities chan ActivityLog
}

func NewActivityLogger(bufferSize int) *ActivityLogger {
	al := &ActivityLogger{
		activities: make(chan ActivityLog, bufferSize),
	}
	go al.processLogs()
	return al
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		activity := ActivityLog{
			Timestamp: time.Now(),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IP:        r.RemoteAddr,
			Duration:  duration,
		}
		
		select {
		case al.activities <- activity:
		default:
			log.Println("Activity log buffer full, dropping entry")
		}
	})
}

func (al *ActivityLogger) processLogs() {
	for activity := range al.activities {
		log.Printf("ACTIVITY: %s %s from %s (%s) took %v",
			activity.Method,
			activity.Path,
			activity.IP,
			activity.UserAgent,
			activity.Duration,
		)
	}
}

func (al *ActivityLogger) Close() {
	close(al.activities)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}