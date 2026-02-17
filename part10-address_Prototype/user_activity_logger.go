package middleware

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

	log.Printf("Activity: %s %s from %s took %v",
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

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		log.Printf("[%s] %s %s %d %v",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
		)
	})
}
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
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activity := ActivityLog{
			Timestamp: time.Now().UTC(),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
		}

		log.Printf("Activity: %s %s from %s (%s)", 
			activity.Method, 
			activity.Path, 
			activity.IPAddress, 
			activity.UserAgent)

		next.ServeHTTP(w, r)
	})
}
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp  time.Time
	UserID     string
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
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
		userID := extractUserID(r)
		
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		
		activity := ActivityLog{
			Timestamp:  time.Now(),
			UserID:     userID,
			Method:     r.Method,
			Path:       r.URL.Path,
			StatusCode: lrw.statusCode,
			Duration:   time.Since(start),
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
		log.Printf("ACTIVITY: %s | User: %s | %s %s | Status: %d | Duration: %v",
			activity.Timestamp.Format(time.RFC3339),
			activity.UserID,
			activity.Method,
			activity.Path,
			activity.StatusCode,
			activity.Duration,
		)
	}
}

func extractUserID(r *http.Request) string {
	if user := r.Context().Value("userID"); user != nil {
		if id, ok := user.(string); ok {
			return id
		}
	}
	return "anonymous"
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}