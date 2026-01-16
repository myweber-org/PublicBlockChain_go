
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

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		logEntry := ActivityLog{
			UserID:    extractUserID(r),
			Path:      r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: r.RemoteAddr,
		}
		
		log.Printf("Activity: %s %s by %s from %s", 
			logEntry.Method, 
			logEntry.Path, 
			logEntry.UserID, 
			logEntry.IPAddress)
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		log.Printf("Request completed in %v", duration)
	})
}

func extractUserID(r *http.Request) string {
	if user := r.Header.Get("X-User-ID"); user != "" {
		return user
	}
	return "anonymous"
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
	
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
	al.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}